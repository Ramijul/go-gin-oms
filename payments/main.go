package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Ramijul/go-gin-oms/payments/rabbitmq"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PAYMENT_STATUS_PAID   = "PAID"
	PAYMENT_STATUS_FAILED = "FAILED"
)

/*
Process the payment
*/
func processPayment(amount float64) string {
	if amount > 1000 {
		return PAYMENT_STATUS_FAILED
	}
	return PAYMENT_STATUS_PAID
}

func init() {
	// fails when env is passed from docker-compose
	err := godotenv.Load()
	if err != nil {
		// test for an env
		if len(os.Getenv("RABBITMQ_CONN_STRING")) == 0 {
			panic("Error loading .env file")
		}
	}

	rabbitmq.IntiallizeVariables()
}

func main() {
	// initialize rabbitmq service for sending response back
	senderConn, senderCh, senderQ := rabbitmq.InitializeRabbitMQ(rabbitmq.RESPONSE_QUEUE)
	msgSenderService := &rabbitmq.RabbitMQService{
		Conn: senderConn,
		Ch:   senderCh,
		Q:    senderQ,
	}
	defer msgSenderService.CloseConnection()

	//consumer for payment request
	consumerConn, consumerCh, consumerQ := rabbitmq.InitializeRabbitMQ(rabbitmq.REQUEST_QUEUE)
	defer consumerConn.Close()
	defer consumerCh.Close()

	msgs, err := consumerCh.Consume(
		consumerQ.Name, // queue
		"payments",     // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	rabbitmq.FailOnError(err, "Failed to register the consumer")
	var forever chan struct{}

	go func() {
		for d := range msgs {
			var paymentInfo rabbitmq.OrderCreateEvent
			// unmarshal the message data
			err := json.Unmarshal(d.Body, &paymentInfo)

			if err != nil {
				log.Print("failed to unmarshal", err)
				rabbitmq.NackMessage(d)
				return
			}

			log.Print("received ", paymentInfo)

			// process the payment
			paymentProcessResult := processPayment(paymentInfo.TotalPrice)

			paymentConfEvent := &rabbitmq.PaymentProcessEvent{
				OrderID:       paymentInfo.OrderID,
				PaymentStatus: paymentProcessResult,
			}

			body, err := json.Marshal(paymentConfEvent)

			if err != nil {
				log.Print("failed to marshal ", paymentConfEvent, err)
				rabbitmq.NackMessage(d)
				return
			}

			//send the result back to orders service
			err = msgSenderService.SendMessage(amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})

			// nack on reply error
			if err != nil {
				log.Print("nacking: ", paymentInfo, err)
				rabbitmq.NackMessage(d)
				return
			}

			//ack the message
			rabbitmq.AckMessage(d)
		}

	}()

	log.Print(" [*] Waiting for messages")

	<-forever
}
