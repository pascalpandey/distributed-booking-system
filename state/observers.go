package state

import (
	"sc4051-server/client"
	"github.com/google/uuid"
)

type Observers map[uuid.UUID]*client.Client // Observers is an array of registered clients

// Notfiies the specified observers
func (observers Observers) Notify(message string) {
	for _, observer := range observers {
		observer.SendMessage(message)
	}
}
