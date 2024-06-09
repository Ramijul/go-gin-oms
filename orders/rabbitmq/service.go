package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	Q    amqp.Queue
}

func (r *RabbitMQService) SendMessage(message amqp.Publishing) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Ch.PublishWithContext(ctx, "", r.Q.Name, false, false, message)
}

func (r *RabbitMQService) CloseConnection() {
	defer r.Conn.Close()
	defer r.Ch.Close()
}
