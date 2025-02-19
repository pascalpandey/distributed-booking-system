package state

import (
	"github.com/distributed-systems-be/client"
	"github.com/google/uuid"
)

type Observers map[uuid.UUID]*client.Client

func (observers Observers) Notify(message string) {
	for _, observer := range observers {
		observer.SendMessage(message)
	}
}
