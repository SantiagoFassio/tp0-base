package common

import (
	"fmt"
	"strings"
)

func parseCSVLine(line string) Bet {
	parts := strings.Split(line, "|")

	return Bet{
		Agency:    parts[0],
		Name:      parts[1],
		Surname:   parts[2],
		Document:  parts[3],
		Birthdate: parts[4],
		Number:    parts[5],
	}
}
