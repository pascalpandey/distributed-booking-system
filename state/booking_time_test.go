package state

import (
	"testing"
)

func TestToMinute(t *testing.T) {
	bookingTime := BookingTime{Day: Monday, Hour: 9, Minute: 30}
	expected := 9*60 + 30
	result := bookingTime.ToMinute()

	if result != expected {
		t.Errorf("For bookingTime %+v, expected %v, but got %v", bookingTime, expected, result)
	}
}

func TestAdd(t *testing.T) {
	bookingTime := BookingTime{Day: Monday, Hour: 9, Minute: 30}
	addTime := BookingTime{Day: 1, Hour: 2, Minute: 45}

	expected := BookingTime{Day: 1, Hour: 12, Minute: 15}
	result := bookingTime.Add(addTime)

	if result != expected {
		t.Errorf("For addTime %+v on bookingTime %+v, expected %+v, but got %+v", addTime, bookingTime, expected, result)
	}
}

func TestSubtract(t *testing.T) {
	bookingTime := BookingTime{Day: Tuesday, Hour: 10, Minute: 30}
	subtractTime := BookingTime{Day: -1, Hour: -5, Minute: -45}

	expected := BookingTime{Day: Monday, Hour: 4, Minute: 45}
	result := bookingTime.Subtract(subtractTime)

	if result != expected {
		t.Errorf("For subtractTime %+v on bookingTime %+v, expected %+v, but got %+v", subtractTime, bookingTime, expected, result)
	}
}
