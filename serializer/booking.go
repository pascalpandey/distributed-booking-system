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

// Formats booking time as a string without days
func formatReadableHourMinute(bookingTime state.BookingTime) string {
	return fmt.Sprintf("%d hour(s) and %d minute(s)", bookingTime.Hour, bookingTime.Minute)
}

// Formats booking time as a string with days
func formatBookingTimeWithDay(bookingTime state.BookingTime) string {
	return fmt.Sprintf("%s %02d:%02d", DayToString(bookingTime.Day), bookingTime.Hour, bookingTime.Minute)
}