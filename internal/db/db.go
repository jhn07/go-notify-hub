package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var DB *sql.DB

func Connect(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	DB = db
	log.Println("🔌 Connected to PostgreSQL")
	return nil
}

// CreateTables creates database tables - should be called only by API server
func CreateTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		message TEXT NOT NULL,
		channels TEXT[] NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("✅ Tables created successfully")
	return nil
}
