package state

import (
	"time"

	"sc4051-server/client"
	"github.com/google/uuid"
)

type Facility = string

type FacilityState struct {
	Bookings  []*Booking
	Observers Observers
}

func (facilityState *FacilityState) QueryAvailability(startTime BookingTime, endTime BookingTime) bool {
	bookings := facilityState.Bookings
	for _, booking := range bookings {
		if booking.intersects(startTime, endTime) {
			return false
		}
	}
	return true
}

func (facilityState *FacilityState) Book(startTime BookingTime, endTime BookingTime) uuid.UUID {
	confirmationId := uuid.New()
	newBooking := Booking{
		StartTime:      startTime,
		EndTime:        endTime,
		ConfirmationId: confirmationId,
	}
	facilityState.Bookings = append(facilityState.Bookings, &newBooking)
	return confirmationId
}

func (facilityState *FacilityState) Cancel(confirmationId uuid.UUID) *Booking {
	for i, booking := range facilityState.Bookings {
		if booking.ConfirmationId == confirmationId {
			facilityState.Bookings = append(facilityState.Bookings[:i], facilityState.Bookings[i+1:]...)
		}
	}
	return nil
}

func (facilityState *FacilityState) RegisterObserver(client *client.Client, duration time.Duration) {
	observerId := uuid.New()
	facilityState.Observers[observerId] = client
	go func() {
		time.Sleep(duration)
		delete(facilityState.Observers, observerId)
	}()
}

func (facilityState *FacilityState) NotifyObservers(message string) {
	for _, observer := range facilityState.Observers {
		observer.SendMessage(message)
	}
}
