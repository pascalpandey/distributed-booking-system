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
	bookingTime.Minute += addTime.Minute
	return bookingTime
}

func (bookingTime BookingTime) Subtract(subtractTime BookingTime) BookingTime {
	bookingTime.Day += subtractTime.Day
	bookingTime.Hour += subtractTime.Hour
	bookingTime.Minute += subtractTime.Minute
	return bookingTime
}
