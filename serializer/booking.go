package serializer

import (
	"fmt"

	"github.com/distributed-systems-be/state"
)

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

func formatBookingTime(bookingTime state.BookingTime) string {
	if bookingTime.Day == 0 {
		return fmt.Sprintf("%d/%d", bookingTime.Hour, bookingTime.Minute)
	}
	return fmt.Sprintf("%s/%d/%d", DayToString(bookingTime.Day), bookingTime.Hour, bookingTime.Minute)
}
