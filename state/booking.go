package state

type Booking struct {
	StartTime      BookingTime // Start time of the booking
	EndTime        BookingTime // End time of the booking
	ConfirmationId string      // Unique identifier for the booking
}

// Checks if the given time range overlaps with this booking
func (booking *Booking) intersects(startTime BookingTime, endTime BookingTime) bool {
    startTimeMin := startTime.ToMinute()
    endTimeMin := endTime.ToMinute()
    bookingStartTimeMin := booking.StartTime.ToMinute()
    bookingEndTimeMin := booking.EndTime.ToMinute()
    
    weekMinutes := 7 * 24 * 60
    
    var queryIntervals [][]int
    if startTimeMin > endTimeMin {
        queryIntervals = append(queryIntervals, []int{startTimeMin, weekMinutes})
        queryIntervals = append(queryIntervals, []int{0, endTimeMin})
    } else {
        queryIntervals = append(queryIntervals, []int{startTimeMin, endTimeMin})
    }
    
    var bookingIntervals [][]int
    if bookingStartTimeMin > bookingEndTimeMin {
        bookingIntervals = append(bookingIntervals, []int{bookingStartTimeMin, weekMinutes})
        bookingIntervals = append(bookingIntervals, []int{0, bookingEndTimeMin})
    } else {
        bookingIntervals = append(bookingIntervals, []int{bookingStartTimeMin, bookingEndTimeMin})
    }
    
    for _, qInterval := range queryIntervals {
        qStart, qEnd := qInterval[0], qInterval[1]
        for _, bInterval := range bookingIntervals {
            bStart, bEnd := bInterval[0], bInterval[1]
            if max(qStart, bStart) < min(qEnd, bEnd) {
                return true
            }
        }
    }
    
    return false
}

// Shifts the booking's start and end time by the given offset
func (booking *Booking) Offset(offsetTime BookingTime) {
    if offsetTime.ToMinute() > 0 {
	    booking.StartTime = booking.StartTime.Add(offsetTime)
	    booking.EndTime = booking.EndTime.Add(offsetTime)
    } else {
        booking.StartTime = booking.StartTime.Subtract(offsetTime)
	    booking.EndTime = booking.EndTime.Subtract(offsetTime)
    }
}

// Extends the booking's end time by the given duration
func (booking *Booking) Extend(extendTime BookingTime) {
	booking.EndTime = booking.EndTime.Add(extendTime)
}
