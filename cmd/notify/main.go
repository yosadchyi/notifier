package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yosadchyi/notifier/pkg/notification"

	"github.com/yosadchyi/notifier/pkg/processing"

	"github.com/yosadchyi/notifier/pkg/iohelper"
)

const (
	intervalHelp = "-i, --interval=5s Notification interval"
)

var url = flag.String("url", "", "URL to post messages to, required")
var intervalStr = flag.String("interval", "5s", intervalHelp)

func init() {
	flag.StringVar(intervalStr, "i", "5s", intervalHelp)
}

func main() {
	flag.Parse()

	if url == nil || *url == "" {
		log.Fatalf("URL is required, please provide valid value")
	}
	interval, err := time.ParseDuration(*intervalStr)
	if err != nil {
		log.Fatalf("Bad value for interval %s, please specify valid duration", *intervalStr)
	}

	reader := os.Stdin
	sigint := make(chan os.Signal)
	messages := make(chan string)
	batches := make(chan string)
	batchedMessages := make(chan string)
	cancel := make(chan struct{}, 1)
	discardBatch := make(chan struct{})
	cancelBatchedRead := make(chan struct{})
	signal.Notify(sigint, syscall.SIGINT)
	parallelProcessor := processing.NewParallelProcessor(processing.ParallelConfig{
		WorkerCount:     2,
		InputBufferSize: 128,
	})
	httpNotificationFunc := notification.HttpFunc(*url)
	parallelProcessor.Start(func(msg string) bool {
		httpNotificationFunc(msg)
		return true
	})
	batcher := processing.NewFileBatcher()
	wg := sync.WaitGroup{}
	timer := time.NewTimer(interval)
	wg.Add(1)

	go func() {
		defer close(messages)

		ok := iohelper.ReadAllLines(reader, func(msg string) bool {
			select {
			case messages <- msg:
				return true
			case <-cancel:
				return false
			}
		})
		if !ok {
			timer.Stop()
			// no need to drain timer channel
			discardBatch <- struct{}{}
		}
	}()

	go func() {
		defer close(batches)

		for {
			select {
			case <-discardBatch:
				batcher.DiscardBatch()
				return
			case msg, ok := <-messages:
				if !ok {
					batches <- batcher.FinalizeBatch()
					return
				}
				batcher.AddToBatch(msg)
			case <-timer.C:
				// there is posibility that batch will be finalized in case of cancellation
				// it still will be just deleted in goroutine below
				batches <- batcher.FinalizeBatch()
				timer.Reset(interval)
			}
		}
	}()

	go func() {
		defer close(batchedMessages)

		for batch := range batches {
			ok := batcher.ReadBatch(batch, func(msg string) bool {
				select {
				case batchedMessages <- msg:
					return true
				case <-cancelBatchedRead:
					return false
				}
			})
			if !ok {
				break
			}
		}
	}()

	go func() {
		defer wg.Done()

		for msg := range batchedMessages {
			parallelProcessor.Enqueue(msg)
		}
	}()

	go func() {
		<-sigint
		log.Println("signal received, shutdown...")
		cancel <- struct{}{}
		cancelBatchedRead <- struct{}{}
	}()

	wg.Wait()
	parallelProcessor.Stop()
	parallelProcessor.Wait()
	log.Println("done")
}
