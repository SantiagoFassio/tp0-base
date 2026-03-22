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
func NewClient(config ClientConfig, bet Bet) *Client {
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

func (c *Client) sendBatch(bets []Bet) {
	// Create the connection to the server in every loop iteration.
	c.createClientSocket()
	defer c.conn.Close()

	msg := SerializeBatch(bets)
	
	// short write safe
	totalSent := 0
	data := []byte(msg)

	// Enviamos el batch al servidor
	for totalSent < len(data) {
		n, err := c.conn.Write(data[totalSent:])
		if err != nil {
			log.Errorf("action: batch_enviada | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		totalSent += n
	}

	// Esperamos la respuesta del servidor después de enviar el batch
	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
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
		batch, err := reader.NextBatch()

		// If EOF is reached or there are no bets in the batch, we finish sending bets
		if err == io.EOF || len(batch) == 0 {
			break
		}

		// Send the batch of bets to the server
		c.sendBatch(batch)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
