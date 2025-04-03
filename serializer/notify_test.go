package serializer

import (
	"testing"
	"sc4051-server/state"
)

func TestNotifyBook(t *testing.T) {
	confirmationId := "12345"
	startTime := state.BookingTime{Hour: 14, Minute: 30}
	endTime := state.BookingTime{Hour: 15, Minute: 30}

	expected := "New booking with confirmation ID 12345 from Monday/14/30 to Monday/15/30"
	result := NotifyBook(confirmationId, startTime, endTime)

	if result != expected {
		t.Errorf("For params %s, %+v, %+v, expected %v, but got %v", confirmationId, startTime, endTime, expected, result)
	}
}

func TestNotifyOffset(t *testing.T) {
	confirmationId := "12345"
	offsetTime := state.BookingTime{Hour: 10, Minute: 0}
	booking := &state.Booking{
		StartTime:      state.BookingTime{Day: state.Monday, Hour: 9, Minute: 0},
		EndTime:        state.BookingTime{Day: state.Monday, Hour: 10, Minute: 0},
		ConfirmationId: "12345",
	}

	expected := "Booking with confirmation ID 12345 was shifted by 10 hour(s) and 0 minute(s), new booking time is from Monday/9/0 to Monday/10/0"
	result := NotifyOffset(confirmationId, booking, offsetTime)

	if result != expected {
		t.Errorf("For params %s, %+v, %+v, expected %v, but got %v", confirmationId, booking, offsetTime, expected, result)
	}
}

func TestNotifyExtend(t *testing.T) {
	confirmationId := "12345"
	extendTime := state.BookingTime{Hour: 16, Minute: 0}
	booking := &state.Booking{
		StartTime:      state.BookingTime{Day: state.Monday, Hour: 9, Minute: 0},
		EndTime:        state.BookingTime{Day: state.Monday, Hour: 10, Minute: 0},
		ConfirmationId: "12345",
	}

	expected := "Booking with confirmation ID 12345 was extended by 16 hour(s) and 0 minute(s), new booking time is from Monday/9/0 to Monday/10/0"
	result := NotifyExtend(confirmationId, booking, extendTime)

	if result != expected {
		t.Errorf("For params %s, %+v, %+v, expected %v, but got %v", confirmationId, booking, extendTime, expected, result)
	}
}

func TestNotifyCancel(t *testing.T) {
	confirmationId := "12345"

	expected := "Booking with confirmation ID 12345 was cancelled"
	result := NotifyCancel(confirmationId)

	if result != expected {
		t.Errorf("For confirmationId %s, expected %v, but got %v", confirmationId, expected, result)
	}
}
