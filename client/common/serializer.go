package common

import (
	"fmt"
	"strings"
)

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

func SerializeBatch(bets []Bet) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%d\n", len(bets)))

	for _, bet := range bets {
		builder.WriteString(SerializeBet(bet))
		builder.WriteString("\n")
	}
	return builder.String()
}
