package handler

import (
	"strings"

	"github.com/distributed-systems-be/deserializer"
)

const (
	QueryAvailability   = "QUERY"
	Book                = "BOOK"
	ShiftBooking        = "SHIFT"
	MonitorAvailability = "MONITOR"
	ExtendBooking       = "EXTEND"
	CancelBooking       = "CANCEL"
)

func HandleMessage(message string) {
	body := strings.Split(message, ",")
	operation := body[1]
	switch operation {
		case QueryAvailability: {
			deserializer.QueryAvailabilty(body[2:])
		} 
		case Book: {

		}
		case ShiftBooking: {
			
		}
		case MonitorAvailability: {

		} 
		case ExtendBooking: {

		}
		case CancelBooking: {
			
		}
	}
}