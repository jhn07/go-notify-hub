package main

import (
	"encoding/json"
	"log"
	"notify-hub/internal/channels"
	"notify-hub/internal/models"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	url := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	// Connect to RabbitMQ
	conn, err := connectRabbitMQ(url)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare queue
	q, err := declareQueue(ch)
	if err != nil {
		log.Fatalf("❌ Queue declare failed: %v", err)
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("❌ Failed to register consumer: %v", err)
	}

	log.Println("👂 Worker is waiting for messages...")

	// Setup graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start processing messages
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Println("⚠️ Channel closed, stopping worker")
				return
			}
			processMessage(msg)

		case <-shutdown:
			log.Println("🛑 Shutdown signal received, stopping worker...")
			return
		}
	}
}

// connectRabbitMQ establishes connection to RabbitMQ
func connectRabbitMQ(url string) (*amqp091.Connection, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Connected to RabbitMQ")
	return conn, nil
}

// declareQueue declares the notification queue
func declareQueue(ch *amqp091.Channel) (amqp091.Queue, error) {
	queue, err := ch.QueueDeclare(
		"notifyhub_queue",
		true,  // durable - queue survives broker restart
		false, // autoDelete - queue is not deleted when last consumer unsubscribes
		false, // exclusive - queue can be accessed by other connections
		false, // noWait - don't wait for server confirmation
		nil,   // arguments
	)
	if err != nil {
		return amqp091.Queue{}, err
	}

	log.Println("✅ Queue declared:", queue.Name)
	return queue, nil
}

// processMessage handles a single message from the queue
func processMessage(msg amqp091.Delivery) {
	var payload models.NotificationRequest

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Println("❌ Failed to parse message:", err)
		return
	}

	log.Printf("📨 Processing notification for user %s via channels: %v\n", payload.UserID, payload.Channels)

	// Send notifications through each channel
	for _, channelName := range payload.Channels {
		notifChannel, err := channels.NewChannel(channelName)
		if err != nil {
			log.Printf("❌ Unsupported channel %s: %v\n", channelName, err)
			continue
		}

		// Send asynchronously
		go func(ch channels.Channel, name string) {
			if err := ch.Send(payload.UserID, payload.Message); err != nil {
				log.Printf("❌ Failed to send via %s: %v\n", name, err)
			}
		}(notifChannel, channelName)
	}
}

// getEnv retrieves environment variables with default values
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
