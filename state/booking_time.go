package state

type BookingTime struct {
	Day    Day // Represents the day of the week
	Hour   int // Represents the hour in 24-hour format
	Minute int // Represents the minute within the hour
}

type Day = int

const (
	Monday Day = iota // Represents days as enum of ints starting from 0
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// Converts the booking time into total minutes since the start of the week
func (bookingTime *BookingTime) ToMinute() int {
	return bookingTime.Day*24*60 + bookingTime.Hour*60 + bookingTime.Minute
}

// Adds the given duration to the current booking time, addTime is positive
func (bookingTime BookingTime) Add(addTime BookingTime) BookingTime {
	bookingTime.Day += addTime.Day
	bookingTime.Hour += addTime.Hour
	if bookingTime.Hour >= 24 {
		bookingTime.Hour %= 24
		bookingTime.Day += 1
	}
	bookingTime.Day = mod(bookingTime.Day, 7)
	bookingTime.Minute += addTime.Minute
	if bookingTime.Minute >= 60 {
		bookingTime.Minute %= 60
		bookingTime.Hour += 1
	}
	return bookingTime
}

// Subtracts the given duration from the current booking time, subtractTime is negative
func (bookingTime BookingTime) Subtract(subtractTime BookingTime) BookingTime {
	bookingTime.Day += subtractTime.Day
	bookingTime.Hour += subtractTime.Hour
	if bookingTime.Hour < 0 {
		bookingTime.Hour += 24
		bookingTime.Day -= 1
	}
	bookingTime.Day = mod(bookingTime.Day, 7)
	bookingTime.Minute += subtractTime.Minute
	if bookingTime.Minute < 0 {
		bookingTime.Minute += 60
		bookingTime.Hour -= 1
	}
	return bookingTime
}

// % function that wraps around negative integers to positive
func mod(a, b int) int {
	return (a%b + b) % b
}