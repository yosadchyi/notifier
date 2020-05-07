package iohelper

import (
	"bufio"
	"io"
	"log"
)

func ReadAllLines(rd io.Reader, handler func(line string) bool) bool {
	br := bufio.NewReader(rd)

	for {
		if line, err := br.ReadString('\n'); err != nil {
			if err == io.EOF {
				return true
			}
			log.Fatalf("reading error: %s", err)
		} else {
			if !handler(line) {
				return false
			}
		}
	}
}
