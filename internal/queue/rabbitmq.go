package queue

import (
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp091.Channel
	queue   amqp091.Queue
}

type NotificationMessage struct {
	ID       string   `json:"id"`
	UserID   string   `json:"user_id"`
	Message  string   `json:"message"`
	Channels []string `json:"channels"`
}

// ConnectPublisher creates a publisher and queue.
func ConnectPublisher(url string) (*Publisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"notifyhub_queue",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	log.Println("📬 Connected to RabbitMQ queue:", q.Name)

	return &Publisher{
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Publisher) Publish(msg NotificationMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
