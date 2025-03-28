package state

type Booking struct {
	StartTime      BookingTime  // Start time of the booking
	EndTime        BookingTime  // End time of the booking
	ConfirmationId string       // Unique identifier for the booking
}

// Checks if the given time range overlaps with this booking
func (booking *Booking) intersects(startTime BookingTime, endTime BookingTime) bool {
	startTimeMin := startTime.ToMinute()
	endTimeMin := endTime.ToMinute()
	bookingStartTimeMin := booking.StartTime.ToMinute()
	bookingEndTimeMin := booking.EndTime.ToMinute()
	return max(startTimeMin, bookingStartTimeMin) < min(endTimeMin, bookingEndTimeMin)
}

// Shifts the booking's start and end time by the given offset
func (booking *Booking) Offset(offsetTime BookingTime) {
	booking.StartTime = booking.StartTime.Add(offsetTime)
	booking.EndTime = booking.EndTime.Add(offsetTime)
}

// Extends the booking's end time by the given duration
func (booking *Booking) Extend(extendTime BookingTime) {
	booking.EndTime = booking.EndTime.Add(extendTime)
}
