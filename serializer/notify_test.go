package serializer

import (
	"testing"
	"sc4051-server/state"
)

func TestNotifyBook(t *testing.T) {
	confirmationId := "12345"
	startTime := state.BookingTime{Hour: 14, Minute: 30}
	endTime := state.BookingTime{Hour: 15, Minute: 30}

	expected := "MONITOR,BOOK,12345,14/30,15/30"
	result := NotifyBook(confirmationId, startTime, endTime)

	if result != expected {
		t.Errorf("For params %s, %+v, %+v, expected %v, but got %v", confirmationId, startTime, endTime, expected, result)
	}
}

func TestNotifyOffset(t *testing.T) {
	confirmationId := "12345"
	offsetTime := state.BookingTime{Hour: 10, Minute: 0}

	expected := "MONITOR,OFFSET,12345,10/0"
	result := NotifyOffset(confirmationId, offsetTime)

	if result != expected {
		t.Errorf("For params %s, %+v, expected %v, but got %v", confirmationId, offsetTime, expected, result)
	}
}

func TestNotifyExtend(t *testing.T) {
	confirmationId := "12345"
	extendTime := state.BookingTime{Hour: 16, Minute: 0}

	expected := "MONITOR,EXTEND,12345,16/0"
	result := NotifyExtend(confirmationId, extendTime)

	if result != expected {
		t.Errorf("For params %s, %+v, expected %v, but got %v", confirmationId, extendTime, expected, result)
	}
}

func TestNotifyCancel(t *testing.T) {
	confirmationId := "12345"

	expected := "MONITOR,CANCEL,12345"
	result := NotifyCancel(confirmationId)

	if result != expected {
		t.Errorf("For confirmationId %s, expected %v, but got %v", confirmationId, expected, result)
	}
}
