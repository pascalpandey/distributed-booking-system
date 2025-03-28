package serializer

import (
	"fmt"

	"sc4051-server/state"
)

// Serializes a notification for booking
func NotifyBook(confirmationId string, startTime state.BookingTime, endTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,BOOK,%s,%s,%s", confirmationId, formatBookingTime(startTime), formatBookingTime(endTime))
}

// Serializes a notification for offsetting a booking
func NotifyOffset(confirmationId string, offsetTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,OFFSET,%s,%s", confirmationId, formatBookingTime(offsetTime))
}

// Serializes a notification for extending a booking
func NotifyExtend(confirmationId string, extendTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,EXTEND,%s,%s", confirmationId, formatBookingTime(extendTime))
}

// Serializes a notification for canceling a booking
func NotifyCancel(confirmationId string) string {
	return fmt.Sprintf("MONITOR,CANCEL,%s", confirmationId)
}
