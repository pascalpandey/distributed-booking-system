package deserializer

import (
	"github.com/distributed-systems-be/state"
)

func QueryAvailabilty(body []string) (state.Room, state.BookingTime, state.BookingTime) {
	return body[2], state.BookingTime{
		
	}, state.BookingTime{
		
	}
}