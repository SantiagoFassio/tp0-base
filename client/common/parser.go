package common

import (
	"fmt"
	"strings"
)

func parseCSVLine(line string) Bet {
	parts := strings.Split(line, "|")

	return Bet{
		Name:      parts[0],
		Surname:   parts[1],
		Document:  parts[2],
		Birthdate: parts[3],
		Number:    parts[4],
	}
}
