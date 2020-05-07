package processing_test

import (
	"testing"

	"github.com/yosadchyi/notifier/pkg/processing"
)

func TestParallelProcessor(t *testing.T) {
	t.Run("WHEN input given AND processed THEN all items collected", func(t *testing.T) {
		expectedSet := map[string]bool{
			"1": false,
			"2": false,
			"3": false,
			"4": false,
		}
		messages := make([]string, 0)
		processor := processing.NewParallelProcessor(processing.ParallelConfig{
			WorkerCount:     2,
			InputBufferSize: 8,
		})
		processor.Start(func(item string) bool {
			messages = append(messages, item)
			return true
		})
		processor.Stop()
		processor.Wait()
		for _, m := range messages {
			v, ok := expectedSet[m]
			if !ok {
				t.Fatalf("Unexpected message received: %s", m)
			}
			if v {
				t.Fatalf("Message received twice: %s", m)
			}
			expectedSet[m] = true
		}
	})
}
