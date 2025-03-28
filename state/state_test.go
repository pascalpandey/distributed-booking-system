package state

import (
	"fmt"
	"testing"
	"time"

	"sc4051-server/client"

	"github.com/google/uuid"
)

func TestQueryAvailabilityState(t *testing.T) {
	state := InitState()
	facility := "TR1"
	startTime := BookingTime{Day: Monday, Hour: 9, Minute: 30}
	endTime := BookingTime{Day: Monday, Hour: 10, Minute: 30}

	result, err := state.QueryAvailability(facility, startTime, endTime)

	if err != nil || result == false {
		t.Errorf("For params %s, %+v, %+v expected (true, nil), but got (%v, %v)", facility, startTime, endTime, result, err)
	}
}

func TestBookState(t *testing.T) {
	state := InitState()
	facility := "TR1"
	startTime := BookingTime{Day: Monday, Hour: 9, Minute: 0}
	endTime := BookingTime{Day: Monday, Hour: 10, Minute: 0}

	_, confirmationId, err := state.Book(facility, startTime, endTime)

	if err != nil || !IsValidConfirmationId(confirmationId) {
		t.Errorf("For params %s, %+v, %+v expected (<some valid confirmationId>, nil), but got (%v, %+v)", facility, startTime, endTime, confirmationId, err)
	}
}

func TestOffsetBooking(t *testing.T) {
	state := InitState()
	confirmationId := "CONF-12345"
	offsetTime := BookingTime{Day: 0, Hour: 1, Minute: 30}

	expected := fmt.Errorf("booking with confirmationId CONF-12345 not found")

	_, err := state.OffsetBooking(confirmationId, offsetTime)

	if err == nil || err.Error() != expected.Error() {
		t.Errorf("For params %v, %+v expected %+v, but got %+v", confirmationId, offsetTime, expected, err)
	}
}

func TestExtendBooking(t *testing.T) {
	state := InitState()
	confirmationId := "CONF-12345"
	extendTime := BookingTime{Day: 0, Hour: 1, Minute: 0}

	expected := fmt.Errorf("booking with confirmationId CONF-12345 not found")

	_, err := state.ExtendBooking(confirmationId, extendTime)

	if err == nil || err.Error() != expected.Error() {
		t.Errorf("For params %v, %+v expected %+v, but got %+v", confirmationId, extendTime, expected, err)
	}
}

func TestCancelBooking(t *testing.T) {
	state := InitState()
	confirmationId := "CONF-12345"

	_, canceled := state.CancelBooking(confirmationId)

	if canceled != true {
		t.Errorf("For confirmationId %+v expected true, but got canceled status %v", confirmationId, canceled)
	}
}

func TestMonitorState(t *testing.T) {
	state := InitState()
	client := &client.Client{}
	facility := "TR1"
	monitorDuration := time.Hour

	err := state.Monitor(client, facility, monitorDuration)

	if err != nil {
		t.Errorf("For params %v, %+v expected (nil), but got (%v)", facility, monitorDuration, err)
	}
}

// Example of a valid confirmationId is CONF-24d2e00c-fada-4e01-84fc-f9ca134ac62c
func IsValidConfirmationId(u string) bool {
	_, err := uuid.Parse(u[5:])
	return err == nil
}
