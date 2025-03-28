package state

import (
	"testing"
)

func TestIntersects(t *testing.T) {
	booking := &Booking{
		StartTime:      BookingTime{Day: Monday, Hour: 9, Minute: 0},
		EndTime:        BookingTime{Day: Monday, Hour: 10, Minute: 0},
		ConfirmationId: "12345",
	}
	startTime := BookingTime{Day: Monday, Hour: 9, Minute: 30}
	endTime := BookingTime{Day: Monday, Hour: 10, Minute: 30}

	expected := true
	result := booking.intersects(startTime, endTime)

	if result != expected {
		t.Errorf("For params %+v, %+v on %+v, expected %v, but got %v", startTime, endTime, *booking, expected, result)
	}
}

func TestOffset(t *testing.T) {
	booking := &Booking{
		StartTime:      BookingTime{Day: Monday, Hour: 9, Minute: 0},
		EndTime:        BookingTime{Day: Monday, Hour: 10, Minute: 0},
		ConfirmationId: "12345",
	}
	offsetTime := BookingTime{Day: 1, Hour: 2, Minute: 30}

	booking.Offset(offsetTime)

	expectedStart := BookingTime{Day: Tuesday, Hour: 11, Minute: 30}
	expectedEnd := BookingTime{Day: Tuesday, Hour: 12, Minute: 30}

	if booking.StartTime != expectedStart || booking.EndTime != expectedEnd {
		t.Errorf("For offsetTime %+v on booking %+v, expected (%+v, %+v), but got (%+v, %+v)", offsetTime, *booking, expectedStart, expectedEnd, booking.StartTime, booking.EndTime)
	}
}

func TestExtend(t *testing.T) {
	booking := &Booking{
		StartTime:      BookingTime{Day: Monday, Hour: 9, Minute: 0},
		EndTime:        BookingTime{Day: Monday, Hour: 10, Minute: 0},
		ConfirmationId: "12345",
	}
	extendTime := BookingTime{Day: 0, Hour: 2, Minute: 30}

	booking.Extend(extendTime)

	expectedEnd := BookingTime{Day: Monday, Hour: 12, Minute: 30}

	if booking.EndTime != expectedEnd {
		t.Errorf("For extendTime %+v on booking %+v, expected %+v, but got %+v", extendTime, *booking, expectedEnd, booking.EndTime)
	}
}
