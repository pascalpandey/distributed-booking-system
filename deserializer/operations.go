package deserializer

import (
	"strconv"
	"time"

	"sc4051-server/state"
)

func FacilityWithBooking(body []string) (state.Facility, state.BookingTime, state.BookingTime) {
	return body[0], deserializeBookingTime(body[1]), deserializeBookingTime(body[2])
}

func ConfirmationIdWithBookingTime(body []string) (string, state.BookingTime) {
	return body[0], deserializeBookingTime(body[1])
}

func ConfirmationId(body []string) string {
	return body[0]
}

func FacilityWithMonitorDuration(body []string) (state.Facility, time.Duration) {
	duration, _ := strconv.Atoi(body[1])
	return body[0], time.Duration(duration) * time.Second
}
