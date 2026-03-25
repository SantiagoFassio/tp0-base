package common

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"os"
	"io"
	"strings"
	"strconv"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	Agency		  string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	MaxAmount     int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// Sends the server the message. The message can be about sending a batch of bets (B)
// or letting the server know the client finished sending bets (E)
func (c *Client) sendAndReceive(msg string) (string, error) {
	c.createClientSocket()
	defer c.conn.Close()

	if err := writeAll(c.conn, []byte(msg)); err != nil {
		return "", err
	}
	
	reader := bufio.NewReader(c.conn)
	return readLine(reader)
}

// Sends the server a message to ask for the bet results. Returns false if an error ocurrs or
// if any client did not finish sending all their bets.
func (c *Client) sendAndReceiveWinners(msg string) ([]string, bool, error) {
	c.createClientSocket()
	defer c.conn.Close()

	if err := writeAll(c.conn, []byte(msg)); err != nil {
		return nil, false, err
	}

	reader := bufio.NewReader(c.conn)
	header, err := readLine(reader)

	if err != nil {
		return nil, false, err
	}

	if header == "NOK" {
		return nil, false, nil
	}

	n, err := strconv.Atoi(header)
	if err != nil {
		return nil, false, err
	}

	winners, err := readNLines(reader, n)
	if err != nil {
		return nil, false, err
	}

	return winners, true, nil
}

// Send batch sends a batch of bets to the client
func (c *Client) sendBatch(bets []Bet) {

	msg := SerializeBatch(bets)
	response, err := c.sendAndReceive(msg)
	
	if err != nil {
		log.Errorf("action: batch_enviada | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	log.Infof("action: receive_message | result: success | client_id: %v | response: %v",
		c.config.ID,
		strings.TrimSpace(response),
	)
}

// Send Bets starts the loop and sends messages to the client until the bets in 
// the file are finished or a shutdown signal is received.
func (c *Client) SendBets(done chan os.Signal) {
	// Path to the file
	path := fmt.Sprintf("/data/agency-%s.csv", c.config.Agency)

	reader, err := NewBatchReader(path, c.config.MaxAmount)
	if err != nil {
		log.Criticalf("action: read_file | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	for {
		// To ensure graceful shutdown when signal is received, we check if the done channel has received a signal
		select {
			case <-done:
				log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
				return
			default:
		}
		
		// Make the reader read the next batch of bets
		batch, err := reader.NextBatch(c.config.Agency)

		// If EOF is reached or there are no bets in the batch, we finish sending bets
		if err == io.EOF || len(batch) == 0 {
			break
		}

		// Send the batch of bets to the server
		c.sendBatch(batch)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	return
}

// SendEnd sends a message that lets the Server know it finished sending all bets.
// Message consists of "E|x\n". E being the signal that marks the type of message, and
// x being the agency number.
func (c *Client) SendEnd(done chan os.Signal) {
	msg := fmt.Sprintf("E|%s\n", c.config.Agency)

	for {
		// To ensure graceful shutdown when signal is received, we check if the done channel has received a signal
		select {
			case <-done:
				log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
				return
			default:
		}

		response, err := c.sendAndReceive(msg)

		if err != nil {
			log.Errorf("action: end_enviado | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		if response == "OK" {
			log.Infof("action: end_enviado | result: success | client_id: %v | agencia: %v",
				c.config.ID, c.config.Agency)
			return
		}

		//retry
		time.Sleep(200 * time.Millisecond)
	}
}

// GetResults makes the client do polling on the server until it receives all winners of the bets the agency sent.
// Messsage consists of "G|x\n". G being the signal that marks the type of message, and
// x being the agency number.
func (c *Client) GetResults(done chan os.Signal) {
	msg := fmt.Sprintf("G|%s\n", c.config.Agency)

	for {
		// To ensure graceful shutdown when signal is received, we check if the done channel has received a signal
		select {
			case <-done:
				log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
				return
			default:
		}

		response, ready, err := c.sendAndReceiveWinners(msg)

		if err != nil {
			log.Errorf("action: consulta_ganadores | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		if !ready {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if err != nil {
			log.Errorf("action: consulta_ganadores | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v",
			len(response))
		return
	}
}
