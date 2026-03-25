package common

import (
	"os"
	"bufio"
	"io"
	"strings"
	"fmt"
)

// BatchReader is responsible for reading bets from a CSV file in batches,
// ensuring that the total size of the batch does not exceed 8KB and the number of bets does not exceed maxAmount.
type BatchReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	pending   *string
	maxAmount int
	maxBytes  int
}

// NewBatchReader initializes a new BatchReader for the given file path and maximum amount of bets per batch.
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

// NextBatch reads the next batch of bets from the CSV file, ensuring that the total size 
// of the batch does not exceed 8KB and the number of bets does not exceed maxAmount.
// It returns a slice of Bet and an error if any occurs during reading.
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
		headerSize := len(fmt.Sprintf("B|%d\n", len(batch)+1)) // +1 for the new bet

		if len(batch) > 0 && (len(batch) >= r.maxAmount || currentSize+lineSize+headerSize > r.maxBytes) {
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
