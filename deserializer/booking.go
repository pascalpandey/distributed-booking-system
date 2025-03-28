package deserializer

import (
	"strconv"
	"strings"

	"sc4051-server/state"
)

// Deserializes a string representing the day of the week to the corresponding Day enum
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

// Deserializes a string in the format "Day/Hour/Minute" or "Hour/Minute" to a booking time type
func deserializeBookingTime(str string) state.BookingTime {
	lst := strings.Split(str, "/")
	if len(lst) == 3 {
		day := stringToDay(lst[0])
		hour, _ := strconv.Atoi(lst[1])
		minute, _ := strconv.Atoi(lst[2])
		return state.BookingTime{
			Day:    day,
			Hour:   hour,
			Minute: minute,
		}
	} else if len(lst) == 2 {
		hour, _ := strconv.Atoi(lst[0])
		minute, _ := strconv.Atoi(lst[1])
		return state.BookingTime{
			Hour:   hour,
			Minute: minute,
		}
	}
	return state.BookingTime{}
}
