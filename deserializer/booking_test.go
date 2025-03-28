package deserializer

import (
	"testing"

	"sc4051-server/state"
)

func TestStringToDay(t *testing.T) {
	tests := []struct {
		input    string
		expected state.Day
	}{
		{"monday", state.Monday},
		{"tuesday", state.Tuesday},
		{"wednesday", state.Wednesday},
		{"thursday", state.Thursday},
		{"friday", state.Friday},
		{"saturday", state.Saturday},
		{"sunday", state.Sunday},
		{"invalid", -1},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := stringToDay(test.input)
			if result != test.expected {
				t.Errorf("For dayStr %s expected %+v, got %+v", test.input, test.expected, result)
			}
		})
	}
}

func TestDeserializeBookingTime(t *testing.T) {
	tests := []struct {
		input    string
		expected state.BookingTime
	}{
		{"monday/14/30", state.BookingTime{Day: state.Monday, Hour: 14, Minute: 30}},
		{"tuesday/10/15", state.BookingTime{Day: state.Tuesday, Hour: 10, Minute: 15}},
		{"14/30", state.BookingTime{Hour: 14, Minute: 30}},
		{"15/45", state.BookingTime{Hour: 15, Minute: 45}},
		{"invalid", state.BookingTime{}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := deserializeBookingTime(test.input)
			if result != test.expected {
				t.Errorf("For str %s expected %+v, got %+v", test.input, test.expected, result)
			}
		})
	}
}
