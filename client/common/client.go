package common

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	bet    Bet
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet Bet) *Client {
	client := &Client{
		config: config,
		bet:    bet,
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

// SerializeBet Serializes the bet struct to a string format to be sent to the server
func (c *Client) SerializeBet() string {
	msg:= fmt.Sprintf("%s|%s|%s|%s|%s\n",
		c.bet.Nombre,
		c.bet.Apellido,
		c.bet.DNI,
		c.bet.Nacimiento,
		c.bet.Numero
	)
	return msg
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(done chan os.Signal) {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// To ensure graceful shutdown when signal is received, we check if the done channel has received a signal
		select {
		case <-done:
			log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
			return
		default:
		}
		
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		msg := c.SerializeBet()
		data := []byte(msg)
		totalWritten := 0

		// DONE: Send message to the server accounting for short writes.
		for totalWritten < len(data) {
			n, err := c.conn.Write(data[totalWritten:])
			if err != nil {
				log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				c.conn.Close()
				return
			}
			totalWritten += n
		}

		response, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | response: %v",
			c.config.ID,
			response,
		)

		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			c.bet.DNI,
			c.bet.Numero,
		)

		// Check for graceful shutdown signal
		// Wait a time between sending one message and the next one
		select {
		case <-done:
			log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
			return
		case <-time.After(c.config.LoopPeriod):
		}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
