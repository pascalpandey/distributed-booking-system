package deserializer

import (
	"testing"
	"time"

	"sc4051-server/state"
)

func TestFacilityWithBooking(t *testing.T) {
	body := []string{"LAB1", "monday/9/0", "monday/10/0"}
	expectedFacility := "LAB1"
	expectedStartTime := state.BookingTime{Day: state.Monday, Hour: 9, Minute: 0}
	expectedEndTime := state.BookingTime{Day: state.Monday, Hour: 10, Minute: 0}

	facility, startTime, endTime := FacilityWithBooking(body)
	if facility != expectedFacility || startTime != expectedStartTime || endTime != expectedEndTime {
		t.Errorf("For body %+v, expected (%v, %+v, %+v), but got (%v, %+v, %+v)",
			body, expectedFacility, expectedStartTime, expectedEndTime, facility, startTime, endTime)
	}
}

func TestConfirmationIdWithBookingTime(t *testing.T) {
	body := []string{"12345", "tuesday/14/0"}
	expectedConfirmationId := "12345"
	expectedBookingTime := state.BookingTime{Day: state.Tuesday, Hour: 14, Minute: 0}

	confirmationId, bookingTime := ConfirmationIdWithBookingTime(body)
	if confirmationId != expectedConfirmationId || bookingTime != expectedBookingTime {
		t.Errorf("For body %+v, expected (%v, %+v), but got (%v, %+v)",
			body, expectedConfirmationId, expectedBookingTime, confirmationId, bookingTime)
	}
}

func TestConfirmationId(t *testing.T) {
	body := []string{"12345"}
	expectedConfirmationId := "12345"

	confirmationId := ConfirmationId(body)
	if confirmationId != expectedConfirmationId {
		t.Errorf("For body %+v, expected %v, but got %v", body, expectedConfirmationId, confirmationId)
	}
}

func TestFacilityWithMonitorDuration(t *testing.T) {
	body := []string{"LAB1", "30"}
	expectedFacility := "LAB1"
	expectedDuration := 30 * time.Second

	facility, duration := FacilityWithMonitorDuration(body)
	if facility != expectedFacility || duration != expectedDuration {
		t.Errorf("For body %+v, expected (%v, %+v), but got (%v, %+v)",
			body, expectedFacility, expectedDuration, facility, duration)
	}
}
