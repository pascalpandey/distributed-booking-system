package serializer

import (
	"fmt"

	"github.com/distributed-systems-be/state"
	"github.com/google/uuid"
)

func NotifyBook(confirmationId uuid.UUID, startTime state.BookingTime, endTime state.BookingTime) string {
	return fmt.Sprintf("NEW,%s,%s,%s", confirmationId.String(), formatBookingTime(startTime), formatBookingTime(endTime))
}

func NotifyOffset(confirmationId uuid.UUID, offsetTime state.BookingTime) string {
	return fmt.Sprintf("OFFSET,%s,%s", confirmationId.String(), formatBookingTime(offsetTime))
}

func NotifyExtend(confirmationId uuid.UUID, extendTime state.BookingTime) string {
	return fmt.Sprintf("EXTEND,%s,%s", confirmationId.String(), formatBookingTime(extendTime))
}

func NotifyCancel(confirmationId uuid.UUID) string {
	return fmt.Sprintf("CANCEL,%s", confirmationId.String())
}