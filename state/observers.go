package state

import (
	"sc4051-server/client"
)

type Observers map[string]*client.Client // Observers is a map of UUID to registered clients

// Notfiies the specified observers
func (observers Observers) Notify(message string) {
	for _, observer := range observers {
		observer.SendMessage(message)
	}
}
