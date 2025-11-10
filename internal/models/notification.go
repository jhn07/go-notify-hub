package models

type NotificationRequest struct {
	UserID   string                 `json:"user_id"`
	Message  string                 `json:"message"`
	Channels []string               `json:"channels"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

type NotificationResponse struct {
	Status    string   `json:"status"`
	MessageID string   `json:"message_id"`
	Channels  []string `json:"channels"`
}
