package iohelper_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yosadchyi/notifier/pkg/iohelper"
)

func TestReadAllLines(t *testing.T) {
	t.Run("All lines read", func(t *testing.T) {
		reader := strings.NewReader("1\n2\n3\n")
		expectedLines := []string{"1\n", "2\n", "3\n"}
		lines := make([]string, 0)
		iohelper.ReadAllLines(reader, func(line string) bool {
			lines = append(lines, line)
			return true
		})
		if !reflect.DeepEqual(expectedLines, lines) {
			t.Fatalf("expected to get %v but got %v", expectedLines, lines)
		}
	})
	t.Run("Read cancelled", func(t *testing.T) {
		reader := strings.NewReader("1\n2\n3\n")
		expectedLines := []string{"1\n", "2\n"}
		lines := make([]string, 0)
		iohelper.ReadAllLines(reader, func(line string) bool {
			lines = append(lines, line)
			return len(lines) < 2
		})
		if !reflect.DeepEqual(expectedLines, lines) {
			t.Fatalf("expected to get %v but got %v", expectedLines, lines)
		}
	})
}
