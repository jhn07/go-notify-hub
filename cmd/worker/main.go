package main

import (
	"encoding/json"
	"log"
	"notify-hub/internal/channels"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	url := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notifyhub_queue",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("❌ Queue declare failed: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Failed to register consumer: %v", err)
	}

	log.Println("👂 Worker is waiting for messages...")

	for msg := range msgs {
		var payload struct {
			UserID   string   `json:"user_id"`
			Message  string   `json:"message"`
			Channels []string `json:"channels"`
		}
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Println("❌ Failed to parse message:", err)
			continue
		}

		for _, chName := range payload.Channels {
			ch, err := channels.NewChannel(chName)
			if err != nil {
				log.Println("❌ Unsupported channel:", chName)
				continue
			}
			go ch.Send(payload.UserID, payload.Message)
		}
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
