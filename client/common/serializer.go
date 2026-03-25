package common

import (
	"fmt"
	"strings"
)

// SerializeBet creates a string from a Bet object.
// Separates each element with "|"
func SerializeBet(bet Bet) string {
	msg:= fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		bet.Agency,
		bet.Name,
		bet.Surname,
		bet.Document,
		bet.Birthdate,
		bet.Number,
	)
	return msg
}

// Serializes an entire batch
// Starts with "B|x\n", B being the signal and x being the amount of bets sent.
// Each bet is serialized and separated from each other using "\n"
func SerializeBatch(bets []Bet) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("B|%d\n", len(bets)))

	for _, bet := range bets {
		builder.WriteString(SerializeBet(bet))
		builder.WriteString("\n")
	}
	return builder.String()
}
