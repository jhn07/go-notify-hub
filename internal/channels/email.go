package channels

import (
	"fmt"
	"time"
)

type EmailChannel struct{}

func (c *EmailChannel) Send(userID, message string) error {
	fmt.Printf("📧 [Email] Sending to %s: %q\n", userID, message)
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("✅ [Email] Delivered to %s\n", userID)
	return nil
}
