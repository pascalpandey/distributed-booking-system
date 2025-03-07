package serializer

import (
	"fmt"

	"sc4051-server/state"
)

func NotifyBook(confirmationId string, startTime state.BookingTime, endTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,BOOK,%s,%s,%s", confirmationId, formatBookingTime(startTime), formatBookingTime(endTime))
}

func NotifyOffset(confirmationId string, offsetTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,OFFSET,%s,%s", confirmationId, formatBookingTime(offsetTime))
}

func NotifyExtend(confirmationId string, extendTime state.BookingTime) string {
	return fmt.Sprintf("MONITOR,EXTEND,%s,%s", confirmationId, formatBookingTime(extendTime))
}

func NotifyCancel(confirmationId string) string {
	return fmt.Sprintf("MONITOR,CANCEL,%s", confirmationId)
}
