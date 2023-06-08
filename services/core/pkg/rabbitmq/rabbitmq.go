package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ(url string) *RabbitMQ {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalln(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("RabbitMQ connected :: %s \n", url)

	return &RabbitMQ{
		Connection: conn,
		channel:    channel,
	}

}

func (m *RabbitMQ) Publish(queueName string, payload []byte) error {
	_, err := m.channel.QueueDeclare(
		queueName, // queue name
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = m.channel.PublishWithContext(
		context.TODO(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
