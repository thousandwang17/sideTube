package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var rabbitClient *RabbitClient

func init() {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RABBIT_USER_NAME"), os.Getenv("RABBIT_USER_PASS"), os.Getenv("RABBIT_HOST"), os.Getenv("RABBIT_PORT"))
	rabbitClient = clientConnect(addr)
}

func GetRabbitClient() *RabbitClient {
	for !rabbitClient.isReady {
		<-time.After(time.Second)
	}
	return rabbitClient
}

type RabbitClient struct {
	logger          *log.Logger
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

const (
	reconnectDelay = 5 * time.Second

	reInitDelay = 2 * time.Second

	resendDelay = 5 * time.Second

	maxWorkerLimit = 32

	maxWorkerHandlerLimit = 10000
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("client is shutting down")
)

// New creates a new consumer state instance, and automatically
// attempts to connect to the server.
func clientConnect(addr string) *RabbitClient {
	client := RabbitClient{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		done:   make(chan bool),
	}
	go client.handleReconnect(addr)
	return &client
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (client *RabbitClient) handleReconnect(addr string) {
	for {
		client.isReady = false
		client.logger.Println("Attempting to connect")

		conn, err := client.connect(addr)

		if err != nil {
			client.logger.Println("Failed to connect. Retrying...")

			select {
			case <-client.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := client.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (client *RabbitClient) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	client.changeConnection(conn)
	client.logger.Println("Connected!")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (client *RabbitClient) handleReInit(conn *amqp.Connection) bool {
	for {
		client.isReady = false

		err := client.init(conn)

		if err != nil {
			client.logger.Println("Failed to initialize channel. Retrying...")

			select {
			case <-client.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-client.done:
			return true
		case <-client.notifyConnClose:
			client.logger.Println("Connection closed. Reconnecting...")
			return false
		case <-client.notifyChanClose:
			client.logger.Println("Channel closed. Re-running init...")
		}
	}
}

// init will initialize channel & declare queue
func (client *RabbitClient) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	if err = insertSearchEngineQueueRegister(ch); err != nil {
		return err
	}

	client.changeChannel(ch)
	client.isReady = true
	client.logger.Println("Setup!")

	return nil
}

func insertSearchEngineQueueRegister(ch *amqp.Channel) error {
	queue := os.Getenv("INSERT_SEARCH_ENGINE_QUEUE")
	if queue == "" {
		return errors.New("dected queue name is empty")
	}

	exchange := os.Getenv("INSERT_SEARCH_ENGINE_EXCHANGE")
	if exchange == "" {
		return errors.New("dected exchange name is empty")
	}

	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		map[string]interface{}{
			// "x-dead-letter-exchange": "my-dlx",
			// "x-dead-letter-routing-key": "my-dlq",
			// "x-message-ttl": 60000,
			"x-max-retries": 3,
		})

	if err != nil {
		log.Println(err, "Failed to declare an dected queue")
		return err
	}

	err = ch.ExchangeDeclare(
		exchange,            // name
		amqp.ExchangeDirect, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		}, // arguments
	)

	if err != nil {
		log.Println(err, "Failed to declare an dected exchange")
		return err
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		q.Name,   // routing key
		exchange, // exchange
		false,
		nil,
	)

	if err != nil {
		log.Println(err, "Failed to Bind dected queue and exchange")
		return err
	}

	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (client *RabbitClient) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (client *RabbitClient) changeChannel(channel *amqp.Channel) {
	client.channel = channel
	client.notifyChanClose = make(chan *amqp.Error, 1)
	client.notifyConfirm = make(chan amqp.Confirmation, 1)
	client.channel.NotifyClose(client.notifyChanClose)
	client.channel.NotifyPublish(client.notifyConfirm)
}

// Consume will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (client *RabbitClient) Consume() (<-chan amqp.Delivery, error) {
	if !client.isReady {
		return nil, errNotConnected
	}

	worker := 1
	if number, err := strconv.Atoi(os.Getenv("NUMBER_OF_WORKER")); err != nil {
		log.Println("NUMBER_OF_WORKER Atoi error :", err)
	} else if number <= maxWorkerLimit && number > 1 {
		worker = number
	}

	worker_handler := 1
	if number, err := strconv.Atoi(os.Getenv("NUMBER_OF_WORKER_HANDLER")); err != nil {
		log.Println("NUMBER_OF_WORKER_HANDLER Atoi error :", err)
	} else if number <= maxWorkerHandlerLimit && number > 1 {
		worker_handler = number
	}

	if err := client.channel.Qos(
		worker*worker_handler, // set number that the server can handle
		0,
		false,
	); err != nil {
		return nil, err
	}

	queue := os.Getenv("INSERT_SEARCH_ENGINE_QUEUE")

	if queue == "" {
		return nil, errors.New("Queue name is empty")
	}

	return client.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

// Close will cleanly shut down the channel and connection.
func (client *RabbitClient) Close() error {
	if !client.isReady {
		return errAlreadyClosed
	}
	close(client.done)
	err := client.channel.Close()
	if err != nil {
		return err
	}
	err = client.connection.Close()
	if err != nil {
		return err
	}

	client.isReady = false
	return nil
}
