package serializer

import (
	"fmt"

	"sc4051-server/state"
)

// Serializes a notification for booking
func NotifyBook(facility, requestId, confirmationId string, startTime state.BookingTime, endTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,BOOK,%s,%s,%s,%s,%s", facility, requestId, confirmationId, formatBookingTimeWithDay(startTime), formatBookingTimeWithDay(endTime))
}

// Serializes a notification for offsetting a booking
func NotifyOffset(confirmationId string, offsetTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,OFFSET,%s,%s", confirmationId, formatBookingTimeWithoutDay(offsetTime))
}

// Serializes a notification for extending a booking
func NotifyExtend(confirmationId string, extendTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,EXTEND,%s,%s", confirmationId, formatBookingTimeWithoutDay(extendTime))
}

// Serializes a notification for canceling a booking
func NotifyCancel(confirmationId string) string {
	return fmt.Sprintf("MONITOR,CANCEL,%s", confirmationId)
}
