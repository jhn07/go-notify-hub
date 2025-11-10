package service

import (
	"log"
	"notify-hub/internal/channels"
)

// NotifyThroughChannels asynchronously sends a message through all specified channels.
func NotifyThroughChannels(channelNames []string, userID, message string) {
	for _, chName := range channelNames {
		ch, err := channels.NewChannel(chName)
		if err != nil {
			log.Printf("❌ Failed to create channel %s: %v", chName, err)
			continue
		}

		go func(c channels.Channel, name string) {
			if err := c.Send(userID, message); err != nil {
				log.Printf("❌ [%s] failed: %v", name, err)
			}
		}(ch, chName)
	}
}
