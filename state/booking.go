package state

import (
	"github.com/google/uuid"
)

type Booking struct {
	StartTime      BookingTime
	EndTime        BookingTime
	confirmationId uuid.UUID
}

func (booking *Booking) intersects(startTime BookingTime, endTime BookingTime) bool {
	startTimeMin := startTime.ToMinute()
	endTimeMin := endTime.ToMinute()
	bookingStartTimeMin := booking.StartTime.ToMinute()
	if bookingStartTimeMin < endTimeMin && bookingStartTimeMin > startTimeMin {
		return true
	}
	return false
}

func (booking *Booking) Offset(offsetTime BookingTime) {
	booking.StartTime = booking.StartTime.Add(offsetTime)
	booking.EndTime = booking.EndTime.Subtract(offsetTime)
}

func (booking *Booking) Extend(extendTime BookingTime) {
	booking.EndTime = booking.EndTime.Subtract(extendTime)
}
