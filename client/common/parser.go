package common

import (
	"strings"
	"strconv"
	"fmt"
)

func parseCSVLine(line string, agency string) Bet {
	parts := strings.Split(line, ",")

	return Bet{
		Agency:    agency,
		Name:      parts[0],
		Surname:   parts[1],
		Document:  parts[2],
		Birthdate: parts[3],
		Number:    parts[4],
	}
}

func parseWinners(line string) ([]string, error) {
	lines := strings.Split(line, "\n")

	if len(lines) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	n, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}

	if len(lines) - 1 < n {
		return nil, fmt.Errorf("invalid winners format")
	}

	return lines[1 : n+1], nil
}
