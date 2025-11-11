package main

import (
	"context"
	"log"
	"net/http"
	"notify-hub/internal/api"
	"notify-hub/internal/db"
	"notify-hub/internal/queue"
	"os"
	"os/signal"
	"time"
)

func main() {
	addr := getEnv("ADDR", ":8080")
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:notifyhub@localhost:5432/notifyhub?sslmode=disable")

	// Connect to PostgreSQL
	if err := db.Connect(dbURL); err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}
	defer db.DB.Close()

	// Create tables (only API creates tables)
	if err := db.CreateTables(); err != nil {
		log.Fatalf("❌ Failed to create tables: %v", err)
	}
	defer db.DB.Close()

	publisher, err := queue.ConnectPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	api.InitPublisher(publisher)

	server := &http.Server{
		Addr:         addr,
		Handler:      api.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 NotifyHub is running on %s\n", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Wait for SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Println("Bye 👋")

}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
