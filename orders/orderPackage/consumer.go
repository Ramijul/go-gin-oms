package orderPackage

import (
	"encoding/json"
	"log"

	"github.com/Ramijul/go-gin-oms/orders/rabbitmq"
)

func ConsumePaymentConfirmation(orderService Service) {
	conn, ch, q := rabbitmq.InitializeRabbitMQ(rabbitmq.RESPONSE_QUEUE)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name,   // queue
		"orders", // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)

	rabbitmq.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// unmarshal the message body
			var paymentConfirmation rabbitmq.PaymentProcessEvent
			err := json.Unmarshal(d.Body, &paymentConfirmation)

			if err != nil {
				log.Print("Failed to unmarshal ", err)
				rabbitmq.NackMessage(d)
				return
			}

			log.Print(paymentConfirmation)
			err = orderService.HandlePaymentConfirmation(paymentConfirmation)
			if err != nil {
				log.Print("Failed to update status ", paymentConfirmation, err)
				rabbitmq.NackMessage(d)
				return
			}

			rabbitmq.AckMessage(d)

		}
	}()

	log.Print(" [*] Waiting for messages")
	<-forever
}
