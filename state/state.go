package state

import (
	"fmt"
	"strconv"
	"time"

	"github.com/distributed-systems-be/client"
	"github.com/google/uuid"
)

type State map[Facility]*FacilityState

func InitState() State {
	state := State{}
	for _, roomType := range []string{"TR", "LAB", "THEATRE"} {
		for i := 1; i <= 10; i++ {
			facilityName := roomType + strconv.Itoa(i)
			state[facilityName] = &FacilityState{
				Bookings: []*Booking{},
				Observers: map[uuid.UUID]*client.Client{},
			}
		}
	}
	return state
}

func (state *State) QueryAvailability(facility Facility, startTime BookingTime, endTime BookingTime) (bool, error) {
	facilityState, found := (*state)[facility]
	if !found {
		return false, fmt.Errorf("facility %v not found", facility)
	}
	return facilityState.QueryAvailability(startTime, endTime), nil
}

func (state *State) Book(facility Facility, startTime BookingTime, endTime BookingTime) (Observers, uuid.UUID, error) {
	facilityState, found := (*state)[facility]
	if !found {
		return nil, uuid.Nil, fmt.Errorf("facility %v not found", facility)
	}
	if !facilityState.QueryAvailability(startTime, endTime) {
		return nil, uuid.Nil, fmt.Errorf("facility %v already booked for that period", facility)
	}

	confirmationId := facilityState.Book(startTime, endTime)
	return facilityState.Observers, confirmationId, nil
}

func (state *State) OffsetBooking(confirmationId uuid.UUID, offsetTime BookingTime) (Observers, error) {
	facilityState, booking := state.getBooking(confirmationId)
	if booking == nil {
		return nil, fmt.Errorf("booking with confirmationId %v not found", confirmationId)
	}

	var reqStart, reqEnd BookingTime
	if offsetTime.ToMinute() > 0 {
		reqStart = booking.EndTime
		reqEnd = booking.EndTime.Add(offsetTime)
	} else {
		reqStart = booking.StartTime.Subtract(offsetTime)
		reqEnd = booking.StartTime
	}
	if !facilityState.QueryAvailability(reqStart, reqEnd) {
		return nil, fmt.Errorf("cannot offset booking with confirmationId %v as there are conflicts", confirmationId)
	}

	booking.Offset(offsetTime)
	return facilityState.Observers, nil
}

func (state *State) ExtendBooking(confirmationId uuid.UUID, extendTime BookingTime) (Observers, error) {
	facilityState, booking := state.getBooking(confirmationId)
	if booking == nil {
		return nil, fmt.Errorf("booking with confirmationId %v not found", confirmationId)
	}

	reqStart := booking.EndTime
	reqEnd := booking.EndTime.Add(extendTime)
	if !facilityState.QueryAvailability(reqStart, reqEnd) {
		return nil, fmt.Errorf("cannot extend booking with confirmationId %v as there are conflicts", confirmationId)
	}

	booking.Extend(extendTime)
	return facilityState.Observers, nil
}

func (state *State) CancelBooking(confirmationId uuid.UUID) (Observers, error) {
	facilityState, booking := state.getBooking(confirmationId)
	if booking == nil {
		return nil, fmt.Errorf("booking with confirmationId %v not found", confirmationId)
	}
	facilityState.Cancel(confirmationId)
	return facilityState.Observers, nil
}

func (state *State) Monitor(client *client.Client, facility Facility, monitorDuration time.Duration) error {
	facilityState, found := (*state)[facility]
	if !found {
		return fmt.Errorf("facility %v not found", facility)
	}
	facilityState.RegisterObserver(client, monitorDuration)
	return nil
}

func (state State) getBooking(confirmationId uuid.UUID) (*FacilityState, *Booking) {
	for facility, facilityState := range state {
		for _, booking := range facilityState.Bookings {
			if booking.confirmationId == confirmationId {
				return state[facility], booking
			}
		}
	}
	return nil, nil
}