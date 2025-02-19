package deserializer

import (
	"time"

	"github.com/distributed-systems-be/state"
	"github.com/google/uuid"
)

func FacilityWithBooking(body []string) (state.Facility, state.BookingTime, state.BookingTime) {
	return body[0], deserializeBookingTime(body[1]), deserializeBookingTime(body[2])
}

func ConfirmationIdWithBookingTime(body []string) (uuid.UUID, state.BookingTime) {
	confirmationId, _ := uuid.Parse(body[0])
	return confirmationId, deserializeBookingTime(body[1])
}

func ConfirmationId(body []string) uuid.UUID {
	confirmationId, _ := uuid.Parse(body[0])
	return confirmationId
}

func FacilityWithMonitorDuration(body []string) (state.Facility, time.Duration) {
	duration, _ := time.ParseDuration(body[1])
	return body[0], duration * time.Second
}
