package deserializer

import (
	"strconv"
	"time"

	"sc4051-server/state"
)

// Deserializes the body of the message to a facility, start time, and end time
func FacilityWithBooking(body []string) (state.Facility, state.BookingTime, state.BookingTime) {
	return body[0], deserializeBookingTime(body[1]), deserializeBookingTime(body[2])
}

// Deserializes the body of the message to a confirmation ID and a booking time
func ConfirmationIdWithBookingTime(body []string) (string, state.BookingTime) {
	return body[0], deserializeBookingTime(body[1])
}

// Deserializes the confirmation ID from the body of the message
func ConfirmationId(body []string) string {
	return body[0]
}

// Deserializes the body of the message to a facility and monitor duration
func FacilityWithMonitorDuration(body []string) (state.Facility, time.Duration) {
	duration, _ := strconv.Atoi(body[1])
	return body[0], time.Duration(duration) * time.Second
}
