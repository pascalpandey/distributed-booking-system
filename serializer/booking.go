package serializer

import (
	"fmt"

	"sc4051-server/state"
)

// Converts a day enum to its string representation
func DayToString(day state.Day) string {
	switch day {
	case state.Monday:
		return "Monday"
	case state.Tuesday:
		return "Tuesday"
	case state.Wednesday:
		return "Wednesday"
	case state.Thursday:
		return "Thursday"
	case state.Friday:
		return "Friday"
	case state.Saturday:
		return "Saturday"
	case state.Sunday:
		return "Sunday"
	default:
		return "Invalid"
	}
}

// Formats booking time as a string
func formatBookingTime(bookingTime state.BookingTime) string {
	if bookingTime.Day == 0 {
		return fmt.Sprintf("%d/%d", bookingTime.Hour, bookingTime.Minute)
	}
	return fmt.Sprintf("%s/%d/%d", DayToString(bookingTime.Day), bookingTime.Hour, bookingTime.Minute)
}
