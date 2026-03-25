package common

import (
	"net"
	"bufio"
	"strings"
)

// WriteAll sends all the data to the connection given to him.
// Accounts for short writes
func writeAll(conn net.Conn, data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		n, err := conn.Write(data[totalSent:])
		if err != nil {
			return err
		}
		totalSent += n
	}
	return nil
}

// readLines reads a line from the reader until a \n is encountered.
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// readNLines reads n lines given by parameters from the reader.
// lines are separated using "\n"
func readNLines(reader *bufio.Reader, n int) ([]string, error) {
	lines := make([]string, 0, n)

	for i := 0; i < n; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		lines = append(lines, strings.TrimSpace(line))
	}

	return lines, nil
}
