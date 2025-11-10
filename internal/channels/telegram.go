package channels

import (
	"fmt"
	"time"
)

type TelegramChannel struct{}

func (c *TelegramChannel) Send(userID, message string) error {
	fmt.Printf("📨 [Telegram] Sending to %s: %q\n", userID, message)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("✅ [Telegram] Delivered to %s\n", userID)
	return nil
}
