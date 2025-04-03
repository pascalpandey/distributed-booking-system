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

// Converts minutes to booking time
func fromMinute(totalMinutes int) BookingTime {
	totalMinutes = mod(totalMinutes, 7*24*60)
	
	day := totalMinutes / (24 * 60)
	totalMinutes %= (24 * 60)
	
	hour := totalMinutes / 60
	minute := totalMinutes % 60
	
	return BookingTime{
	  	Day:    day,
	  	Hour:   hour,
	  	Minute: minute,
	}
}
  
// Adds the given duration to the current booking time
func (bookingTime BookingTime) Add(addTime BookingTime) BookingTime {
	totalMinutes := bookingTime.ToMinute() + addTime.ToMinute()
	return fromMinute(totalMinutes)
}
  
// Subtracts the given duration from the current booking time
func (bookingTime BookingTime) Subtract(subtractTime BookingTime) BookingTime {
	totalMinutes := bookingTime.ToMinute() + subtractTime.ToMinute()
	return fromMinute(totalMinutes)
}

// % function that wraps around negative integers to positive
func mod(a, b int) int {
	return (a%b + b) % b
}