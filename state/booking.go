package state

import (
	"github.com/google/uuid"
)

type Booking struct {
	startTime      BookingTime
	endTime        BookingTime
	confirmationId uuid.UUID
}

type Day int

const (
	Monday Day = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type BookingTime struct {
	Day    Day
	Hour   int
	Minute int
}
