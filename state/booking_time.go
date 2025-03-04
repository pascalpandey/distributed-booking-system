package state

type BookingTime struct {
	Day    Day
	Hour   int
	Minute int
}

type Day = int

const (
	Monday Day = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (bookingTime *BookingTime) ToMinute() int {
	return bookingTime.Day*24*60 + bookingTime.Hour*60 + bookingTime.Minute
}

func (bookingTime BookingTime) Add(addTime BookingTime) BookingTime {
	bookingTime.Day += addTime.Day
	bookingTime.Hour += addTime.Hour
	if bookingTime.Hour >= 24 {
		bookingTime.Hour %= 24
		bookingTime.Day += 1
	}
	bookingTime.Minute += addTime.Minute
	if bookingTime.Minute >= 60 {
		bookingTime.Minute %= 60
		bookingTime.Hour += 1
	}
	return bookingTime
}

func (bookingTime BookingTime) Subtract(subtractTime BookingTime) BookingTime {
	bookingTime.Day += subtractTime.Day
	bookingTime.Hour += subtractTime.Hour
	if bookingTime.Hour < 0 {
		bookingTime.Hour %= 24
		bookingTime.Day -= 1
	}
	bookingTime.Minute += subtractTime.Minute
	if bookingTime.Minute < 60 {
		bookingTime.Minute %= 60
		bookingTime.Hour -= 1
	}
	return bookingTime
}
