package channels

import (
	"fmt"
	"strings"
)

// Channel is a common interface for all notification channels.
type Channel interface {
	Send(userID, message string) error // sends a notification to the user
}

// NewChannel is a factory that creates the appropriate channel by name.
func NewChannel(name string) (Channel, error) {
	switch strings.ToLower(name) {
	case "telegram":
		return &TelegramChannel{}, nil
	case "email":
		return &EmailChannel{}, nil
	}
	return nil, fmt.Errorf("unsupported channel: %s", name)
}
