package serializer

import (
	"fmt"

	"sc4051-server/state"
)

// Serializes a notification for booking
func NotifyBook(confirmationId string, startTime state.BookingTime, endTime state.BookingTime) string {
	return fmt.Sprintf("New booking with confirmation ID %s from %s to %s.", 
		confirmationId, formatBookingTimeWithDay(startTime), formatBookingTimeWithDay(endTime))
}

// Serializes a notification for offsetting a booking
func NotifyOffset(confirmationId string, booking *state.Booking, offsetTime state.BookingTime) string {
	return fmt.Sprintf("Booking with confirmation ID %s was shifted by %s. The new booking time is from %s to %s.",
		confirmationId, formatReadableHourMinute(offsetTime), 
		formatBookingTimeWithDay(booking.StartTime), formatBookingTimeWithDay(booking.EndTime))
}

// Serializes a notification for extending a booking
func NotifyExtend(confirmationId string, booking *state.Booking, extendTime state.BookingTime) string {
	return fmt.Sprintf("Booking with confirmation ID %s was extended by %s. The new booking time is from %s to %s.",
		confirmationId, formatReadableHourMinute(extendTime), 
		formatBookingTimeWithDay(booking.StartTime), formatBookingTimeWithDay(booking.EndTime))
}

// Serializes a notification for canceling a booking
func NotifyCancel(confirmationId string) string {
	return fmt.Sprintf("Booking with confirmation ID %s was cancelled.", confirmationId)
}
