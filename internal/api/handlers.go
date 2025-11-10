package api

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"notify-hub/internal/models"
	"notify-hub/internal/queue"
	"strings"
)

var publisher *queue.Publisher

// InitPublisher initializes the RabbitMQ Publisher
func InitPublisher(p *queue.Publisher) {
	publisher = p
}

// generateID generates a simple UUID-like identifier without external dependencies
func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := cryptoRand.Read(b); err != nil {
		return "", err
	}

	// 32-character hex string
	return "msg_" + hex.EncodeToString(b), nil
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Require JSON content type
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	// Limit request body size (1 MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var req models.NotificationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validateRequest(req); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Generate message ID
	msgID, err := generateID()
	if err != nil {
		http.Error(w, "Failed to generate message ID", http.StatusInternalServerError)
		return
	}

	if publisher == nil {
		http.Error(w, "queue not initialized", http.StatusInternalServerError)
		return
	}

	// Publish to RabbitMQ queue
	err = publisher.Publish(queue.NotificationMessage{
		UserID:   req.UserID,
		Message:  req.Message,
		Channels: req.Channels,
	})
	if err != nil {
		http.Error(w, "failed to publish to queue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := models.NotificationResponse{
		Status:    "queued",
		MessageID: msgID,
		Channels:  req.Channels,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 Accepted - queued for processing
	_ = json.NewEncoder(w).Encode(resp)
}

func validateRequest(req models.NotificationRequest) error {
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("user_id is required")
	}
	if strings.TrimSpace(req.Message) == "" {
		return errors.New("message is required")
	}
	if len(req.Channels) == 0 {
		return errors.New("channels must not be empty")
	}

	// Check for supported channels
	allowed := map[string]bool{
		"telegram": true,
		"email":    true,
	}

	for _, ch := range req.Channels {
		if !allowed[strings.ToLower(ch)] {
			return errors.New("unsupported channel: " + ch)
		}
	}
	return nil
}

// Router creates and configures the HTTP router
func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/send", SendNotificationHandler)
	return LoggingMiddleware(mux)
}
