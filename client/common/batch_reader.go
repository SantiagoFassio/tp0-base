package common

import (
	"os"
	"bufio"
	"io"
	"strings"
)

type BatchReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	pending   *string
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
		pending:   nil,
		maxAmount: maxAmount,
		maxBytes:  8 * 1024, // 8KB
	}, nil
}

func (r *BatchReader) NextBatch(agency string) ([]Bet, error) {
	var batch []Bet
	currentSize := 0

	for {
		var line string
		if r.pending != nil {
			line = *r.pending
			r.pending = nil
		} else {
			if !r.scanner.Scan() {
				break // EOF reached
			}
			line = r.scanner.Text()
		}

		if strings.TrimSpace(line) == "" {
			continue
		}
		
		bet := parseCSVLine(line, agency)

		serialized := SerializeBet(bet)
		lineSize := len(serialized) + 1 // +1 for newline

		if len(batch) > 0 && (len(batch) >= r.maxAmount || currentSize+lineSize > r.maxBytes) {
			auxLine := line
			r.pending = &auxLine
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
