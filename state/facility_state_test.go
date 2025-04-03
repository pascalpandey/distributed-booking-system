package state

import (
	"testing"
	"time"
	"sc4051-server/client"
)

func TestQueryAvailabilityFacilityState(t *testing.T) {
	facilityState := &FacilityState{
		Bookings: []*Booking{
			{
				StartTime:      BookingTime{Day: Monday, Hour: 9, Minute: 0},
				EndTime:        BookingTime{Day: Monday, Hour: 10, Minute: 0},
				ConfirmationId: "CONF-12345",
			},
		},
	}

	startTime := BookingTime{Day: Monday, Hour: 9, Minute: 30}
	endTime := BookingTime{Day: Monday, Hour: 10, Minute: 30}

	expected := false
	result := facilityState.QueryAvailability(startTime, endTime, "")

	if result != expected {
		t.Errorf("For params %+v, %+v, expected %v, but got %v", startTime, endTime, expected, result)
	}
}

func TestBookFacilityState(t *testing.T) {
	facilityState := &FacilityState{}

	facility := "TR1"
	startTime := BookingTime{Day: Monday, Hour: 9, Minute: 0}
	endTime := BookingTime{Day: Monday, Hour: 10, Minute: 0}

	confirmationId := facilityState.Book(facility, startTime, endTime)

	if len(facilityState.Bookings) != 1 {
		t.Errorf("For params %s, %+v, %+v, expected 1 booking, but got %v", facility, startTime, endTime, len(facilityState.Bookings))
	}

	if facilityState.Bookings[0].ConfirmationId != confirmationId {
		t.Errorf("For params %s, %+v, %+v, expected confirmationId %v, but got %v", facility, startTime, endTime, confirmationId, facilityState.Bookings[0].ConfirmationId)
	}
}

func TestCancel(t *testing.T) {
	facilityState := &FacilityState{
		Bookings: []*Booking{
			{
				StartTime:      BookingTime{Day: Monday, Hour: 9, Minute: 0},
				EndTime:        BookingTime{Day: Monday, Hour: 10, Minute: 0},
				ConfirmationId: "CONF-12345",
			},
		},
	}

	confirmationId := "CONF-12345"
	facilityState.Cancel(confirmationId)

	if len(facilityState.Bookings) != 0 {
		t.Errorf("For canceling %v, expected 0 bookings, but got %v", confirmationId, len(facilityState.Bookings))
	}
}

func TestRegisterObserver(t *testing.T) {
	facilityState := &FacilityState{
		Observers: make(Observers),
	}

	client := &client.Client{}
	duration := time.Second * 999

	facilityState.RegisterObserver(client, duration)

	if len(facilityState.Observers) != 1 {
		t.Errorf("For registering %+v , expected 1 observer, but got %v", *client, len(facilityState.Observers))
	}
}
