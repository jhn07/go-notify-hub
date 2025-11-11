package models

type Status string

// Notification status constants
const (
	StatusQueued  Status = "queued"
	StatusSending Status = "sending"
	StatusSent    Status = "sent"
	StatusFailed  Status = "failed"
	StatusPartial Status = "partial"
)

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusQueued, StatusSending, StatusSent, StatusFailed, StatusPartial:
		return true
	}
	return false
}
