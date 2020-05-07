package processing

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/yosadchyi/notifier/pkg/iohelper"
)

type Batcher interface {
	ReadBatch(key string, callback func(msg string) bool) bool
	AddToBatch(msg string)
	FinalizeBatch() string
	DiscardBatch()
}

type fileBatcher struct {
	currentBatchFile *os.File
}

func NewFileBatcher() Batcher {
	batcher := &fileBatcher{}
	batcher.newBatch()
	return batcher
}

func (b *fileBatcher) ReadBatch(key string, callback func(msg string) bool) bool {
	log.Printf("reading batch %s", key)
	file, err := os.Open(key)
	if err != nil {
		log.Fatalf("can't open batch file %s", key)
	}

	ok := iohelper.ReadAllLines(file, callback)

	closeFile(file)
	removeFile(file)
	if !ok {
		log.Printf("batch %s reading aborted", key)
	} else {
		log.Printf("batch %s read", key)
	}
	return ok
}

func (b *fileBatcher) newBatch() {
	log.Printf("new batch")
	if newBatchFile, err := ioutil.TempFile(os.TempDir(), "chunk_*"); err != nil {
		log.Fatalf("can't create new batch: %s", err)
	} else {
		b.currentBatchFile = newBatchFile
	}
}

func (b *fileBatcher) AddToBatch(msg string) {
	if _, err := b.currentBatchFile.WriteString(msg); err != nil {
		log.Fatalf("error adding message to the batch: %s", err)
	}
}

func (b *fileBatcher) FinalizeBatch() string {
	log.Printf("finalize batch")
	closeFile(b.currentBatchFile)
	fileName := b.currentBatchFile.Name()
	b.newBatch()
	return fileName
}

func (b *fileBatcher) DiscardBatch() {
	log.Printf("discard batch")
	closeFile(b.currentBatchFile)
	removeFile(b.currentBatchFile)
}

func removeFile(file *os.File) {
	if err := os.Remove(file.Name()); err != nil {
		log.Fatalf("error removing batch file %s", err)
	}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Fatalf("error closing batch file %s", err)
	}
}
