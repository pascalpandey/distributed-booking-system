package serializer

import (
	"testing"

	"sc4051-server/state"
)

func TestDayToString(t *testing.T) {
	tests := []struct {
		day              state.Day
		expectedString   string
	}{
		{state.Monday, "Monday"},
		{state.Tuesday, "Tuesday"},
		{state.Wednesday, "Wednesday"},
		{state.Thursday, "Thursday"},
		{state.Friday, "Friday"},
		{state.Saturday, "Saturday"},
		{state.Sunday, "Sunday"},
		{state.Day(0), "Invalid"},
	}

	for _, test := range tests {
		result := DayToString(test.day)
		if result != test.expectedString {
			t.Errorf("For day %v, expected %v, but got %v", test.day, test.expectedString, result)
		}
	}
}

func TestFormatBookingTime(t *testing.T) {
	tests := []struct {
		bookingTime       state.BookingTime
		expectedFormatted string
	}{
		{
			bookingTime:       state.BookingTime{Day: state.Monday, Hour: 9, Minute: 30},
			expectedFormatted: "Monday/9/30",
		},
		{
			bookingTime:       state.BookingTime{Hour: 14, Minute: 45},
			expectedFormatted: "14/45",
		},
	}

	for _, test := range tests {
		result := formatBookingTime(test.bookingTime)
		if result != test.expectedFormatted {
			t.Errorf("For bookingTime %+v, expected %v, but got %v", test.bookingTime, test.expectedFormatted, result)
		}
	}
}
