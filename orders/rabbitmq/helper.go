package rabbitmq

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	REQUEST_QUEUE        string
	RESPONSE_QUEUE       string
	RABBITMQ_CONN_STRING string
)

func IntiallizeVariables() {
	REQUEST_QUEUE = os.Getenv("REQUEST_QUEUE")
	RESPONSE_QUEUE = os.Getenv("RESPONSE_QUEUE")
	RABBITMQ_CONN_STRING = os.Getenv("RABBITMQ_CONN_STRING")
}

type OrderCreateEvent struct {
	OrderID    string  `json:"order_id"`
	TotalPrice float64 `json:"total_price"`
}

type PaymentProcessEvent struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"` //SUCCEEDED or FAILED
}

func NackMessage(d amqp.Delivery) {
	err := d.Nack(false, false)

	if err != nil {
		log.Print("nacking failed", err)
	}
}

func AckMessage(d amqp.Delivery) {
	err := d.Ack(false)

	if err != nil {
		log.Print("ack error: ", err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func InitializeRabbitMQ(queueName string) (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	if len(queueName) == 0 {
		panic("Queue name is empty. Please verify it is not missing from env file")
	}

	conn, err := amqp.Dial(RABBITMQ_CONN_STRING)
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	return conn, ch, q
}
