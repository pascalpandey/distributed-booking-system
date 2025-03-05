package state

type Booking struct {
	StartTime      BookingTime
	EndTime        BookingTime
	ConfirmationId string
}

func (booking *Booking) intersects(startTime BookingTime, endTime BookingTime) bool {
	startTimeMin := startTime.ToMinute()
	endTimeMin := endTime.ToMinute()
	bookingStartTimeMin := booking.StartTime.ToMinute()
	bookingEndTimeMin := booking.EndTime.ToMinute()
	if (startTimeMin <= bookingEndTimeMin && startTimeMin >= bookingStartTimeMin) ||
		(endTimeMin <= bookingEndTimeMin && endTimeMin >= bookingStartTimeMin) {
		return true
	}
	return false
}

func (booking *Booking) Offset(offsetTime BookingTime) {
	booking.StartTime = booking.StartTime.Add(offsetTime)
	booking.EndTime = booking.EndTime.Add(offsetTime)
}

func (booking *Booking) Extend(extendTime BookingTime) {
	booking.EndTime = booking.EndTime.Add(extendTime)
}
