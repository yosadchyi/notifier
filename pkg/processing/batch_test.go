package processing_test

import (
	"reflect"
	"testing"

	"github.com/yosadchyi/notifier/pkg/processing"
)

func TestFileBatcher(t *testing.T) {
	t.Run("WHEN all messages accepted AND batch finalized THEN all messages saved", func(t *testing.T) {
		expectedMessages := []string{
			"1\n",
			"2\n",
			"3\n",
		}
		messages := make([]string, 0)
		batcher := processing.NewFileBatcher()
		batcher.AddToBatch("1\n")
		batcher.AddToBatch("2\n")
		batcher.AddToBatch("3\n")
		batch := batcher.FinalizeBatch()
		batcher.ReadBatch(batch, func(msg string) bool {
			messages = append(messages, msg)
			return true
		})
		batcher.DiscardBatch()

		if !reflect.DeepEqual(expectedMessages, messages) {
			t.Fatalf("expected to get %v but got %v", expectedMessages, messages)
		}
	})
}
