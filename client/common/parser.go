package common

import (
	"strings"
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
