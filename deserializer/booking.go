package deserializer

import (
	"strconv"
	"strings"

	"github.com/distributed-systems-be/state"
)

func stringToDay(dayStr string) state.Day {
	switch strings.ToLower(dayStr) {
	case "monday":
		return state.Monday
	case "tuesday":
		return state.Tuesday
	case "wednesday":
		return state.Wednesday
	case "thursday":
		return state.Thursday
	case "friday":
		return state.Friday
	case "saturday":
		return state.Saturday
	case "sunday":
		return state.Sunday
	default:
		return -1
	}
}

func deserializeBookingTime(str string) state.BookingTime {
	lst := strings.Split(str, "/")
	day := stringToDay(lst[0])
	hour, _ := strconv.Atoi(str)
	return state.BookingTime{
		Day: day,
		Hour: hour,

	}
}