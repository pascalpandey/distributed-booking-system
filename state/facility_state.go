package state

import (
	"time"
	"fmt"

	"sc4051-server/client"
	"github.com/google/uuid"
)

type Facility = string

type FacilityState struct {
	Bookings  []*Booking  // Stores all bookings for the facility
	Observers Observers   // Tracks registered observers for the facility
}

// Checks if the given time range is available for booking
func (facilityState *FacilityState) QueryAvailability(startTime BookingTime, endTime BookingTime) bool {
	bookings := facilityState.Bookings
	for _, booking := range bookings {
		if booking.intersects(startTime, endTime) {
			return false
		}
	}
	return true
}

// Creates a new booking for the given time range and returns a confirmation ID
func (facilityState *FacilityState) Book(facility Facility, startTime BookingTime, endTime BookingTime) string {
	confirmationId := fmt.Sprintf("CONF-%s-%s", facility, uuid.New().String())
	newBooking := Booking{
		StartTime:      startTime,
		EndTime:        endTime,
		ConfirmationId: confirmationId,
	}
	facilityState.Bookings = append(facilityState.Bookings, &newBooking)
	return confirmationId
}

// Cancels a booking by confirmation ID and removes it from the list
func (facilityState *FacilityState) Cancel(confirmationId string) *Booking {
	for i, booking := range facilityState.Bookings {
		if booking.ConfirmationId == confirmationId {
			facilityState.Bookings = append(facilityState.Bookings[:i], facilityState.Bookings[i+1:]...)
		}
	}
	return nil
}

// Registers an observer client for a specific duration before being removed automatically
func (facilityState *FacilityState) RegisterObserver(client *client.Client, duration time.Duration) {
	observerId := uuid.New()
	facilityState.Observers[observerId] = client
	go func() {
		time.Sleep(duration)
		delete(facilityState.Observers, observerId)
	}()
}

// Sends a message to all registered observers
func (facilityState *FacilityState) NotifyObservers(message string) {
	for _, observer := range facilityState.Observers {
		observer.SendMessage(message)
	}
}
