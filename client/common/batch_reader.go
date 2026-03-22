package common

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"io"
)

type BatchReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	maxAmount int
	maxBytes  int
}

func NewBatchReader(filePath string, maxAmount int) (*BatchReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &BatchReader{
		file:      file,
		scanner:   bufio.NewScanner(file),
		maxAmount: maxAmount,
		maxBytes:  8 * 1024, // 8KB
	}, nil
}

func (r *BatchReader) NextBatch() ([]Bet, error) {
	var batch []Bet
	currentSize := 0

	for r.scanner.Scan() {
		line := r.scanner.Text()
		
		bet := parseCSVLine(line)

		serialized := SerializeBet(bet)
		lineSize := len(serialized) + 1 // +1 for newline

		if len(batch) >= r.maxAmount || currentSize+lineSize > r.maxBytes {
			return batch, nil
		}
		
		batch = append(batch, bet)
		currentSize += lineSize
	}

	if len(batch) > 0 {
		return batch, nil
	}

	return nil, io.EOF
}
