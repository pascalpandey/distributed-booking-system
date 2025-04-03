package state

import (
	"fmt"
	"strconv"
	"time"

	"sc4051-server/client"
)

type State map[Facility]*FacilityState // Map of facility -> Facility's observers and bookings

// Initializes the state with predefined facility names
func InitState() State {
	state := State{}
	for _, roomType := range []string{"TR", "LAB", "THEATRE"} {
		for i := 1; i <= 10; i++ {
			facilityName := roomType + strconv.Itoa(i)
			state[facilityName] = &FacilityState{
				Bookings:  []*Booking{},
				Observers: map[string]*client.Client{},
			}
		}
	}
	return state
}

// Checks if a facility is available for a given time range
func (state *State) QueryAvailability(facility Facility, startTime BookingTime, endTime BookingTime) (bool, error) {
	facilityState, found := (*state)[facility]
	if !found {
		return false, fmt.Errorf("facility %v not found", facility)
	}
	return facilityState.QueryAvailability(startTime, endTime, ""), nil
}

// Attempts to make a booking for a facility and returns observers and a confirmation ID
func (state *State) Book(facility Facility, startTime BookingTime, endTime BookingTime) (Observers, string, error) {
	facilityState, found := (*state)[facility]
	if !found {
		return nil, "", fmt.Errorf("facility %v not found", facility)
	}
	if !facilityState.QueryAvailability(startTime, endTime, "") {
		return nil, "", fmt.Errorf("facility %v already booked for that period", facility)
	}

	confirmationId := facilityState.Book(facility, startTime, endTime)
	return facilityState.Observers, confirmationId, nil
}

// Shifts a booking's start and end times by the given offset if possible
func (state *State) OffsetBooking(confirmationId string, offsetTime BookingTime) (Observers, *Booking, error) {
	facilityState, booking := state.getBooking(confirmationId)
	if booking == nil {
		return nil, nil, fmt.Errorf("booking with confirmationId %v not found", confirmationId)
	}

	var reqStart, reqEnd BookingTime
	if offsetTime.ToMinute() > 0 {
		reqStart = booking.StartTime.Add(offsetTime)
		reqEnd = booking.EndTime.Add(offsetTime)
	} else {
		reqStart = booking.StartTime.Subtract(offsetTime)
		reqEnd = booking.EndTime.Subtract(offsetTime)
	}
	if !facilityState.QueryAvailability(reqStart, reqEnd, confirmationId) {
		return nil, nil, fmt.Errorf("cannot offset booking with confirmationId %v as there are conflicts", confirmationId)
	}

	booking.Offset(offsetTime)
	return facilityState.Observers, booking, nil
}

// Extends the booking's end time if there is availability
func (state *State) ExtendBooking(confirmationId string, extendTime BookingTime) (Observers, *Booking, error) {
	facilityState, booking := state.getBooking(confirmationId)
	if booking == nil {
		return nil, nil, fmt.Errorf("booking with confirmationId %v not found", confirmationId)
	}

	reqStart := booking.EndTime
	reqEnd := booking.EndTime.Add(extendTime)
	if !facilityState.QueryAvailability(reqStart, reqEnd, confirmationId) {
		return nil, nil, fmt.Errorf("cannot extend booking with confirmationId %v as there are conflicts", confirmationId)
	}

	booking.Extend(extendTime)
	return facilityState.Observers, booking, nil
}

// Removes a booking and notifies observers
func (state *State) CancelBooking(confirmationId string) (Observers, bool) {
	facilityState, booking := state.getBooking(confirmationId)
	if booking == nil {
		return nil, true
	}
	facilityState.Cancel(confirmationId)
	return facilityState.Observers, false
}

// Registers a client to observe changes to a facility's bookings for a duration
func (state *State) Monitor(client *client.Client, facility Facility, monitorDuration time.Duration) error {
	facilityState, found := (*state)[facility]
	if !found {
		return fmt.Errorf("facility %v not found", facility)
	}
	facilityState.RegisterObserver(client, monitorDuration)
	return nil
}

// Searches for a booking by its confirmation ID and returns the facility state and booking
func (state State) getBooking(confirmationId string) (*FacilityState, *Booking) {
	for facility, facilityState := range state {
		for _, booking := range facilityState.Bookings {
			if booking.ConfirmationId == confirmationId {
				return state[facility], booking
			}
		}
	}
	return nil, nil
}
