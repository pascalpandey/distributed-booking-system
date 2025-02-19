package handler

import (
	"github.com/distributed-systems-be/client"
	"github.com/distributed-systems-be/deserializer"
	"github.com/distributed-systems-be/serializer"
	"github.com/distributed-systems-be/state"
)

const (
	QueryAvailability   = "QUERY"
	Book                = "BOOK"
	OffsetBooking       = "OFFSET"
	MonitorAvailability = "MONITOR"
	ExtendBooking       = "EXTEND"
	CancelBooking       = "CANCEL"
)

type Handler struct {
	State         *state.State
	CallingClient *client.Client
	Cache         map[string]string
}

func (handler *Handler) HandleMessage(message string) {
	requestId, operation, body := deserializer.Message(message)

	if handler.Cache != nil {
		reply, requestHandled := handler.Cache[requestId]
		if requestHandled {
			handler.CallingClient.SendMessage(reply)
		}
	}

	switch operation {
	case QueryAvailability:
		facility, startTime, endTime := deserializer.FacilityWithBooking(body)
		available, err := handler.State.QueryAvailability(facility, startTime, endTime)
		reply := serializer.ReplyQueryAvailability(available, err)
		handler.CallingClient.SendMessage(reply)

	case Book:
		facility, startTime, endTime := deserializer.FacilityWithBooking(body)
		observers, confirmationId, err := handler.State.Book(facility, startTime, endTime)
		if err == nil {
			notification := serializer.NotifyBook(confirmationId, startTime, endTime)
			observers.Notify(notification)
		}
		reply := serializer.ReplyBook(confirmationId, err)
		handler.CallingClient.SendMessage(reply)

	case OffsetBooking:
		confirmationId, offsetTime := deserializer.ConfirmationIdWithBookingTime(body)
		observers, err := handler.State.OffsetBooking(confirmationId, offsetTime)
		if err == nil {
			notification := serializer.NotifyOffset(confirmationId, offsetTime)
			observers.Notify(notification)
		}
		reply := serializer.ReplyStatus(err)
		handler.CallingClient.SendMessage(reply)

	case MonitorAvailability:
		facility, monitorDuration := deserializer.FacilityWithMonitorDuration(body)
		err := handler.State.Monitor(handler.CallingClient, facility, monitorDuration)
		reply := serializer.ReplyStatus(err)
		handler.CallingClient.SendMessage(reply)

	case ExtendBooking:
		confirmationId, extendTime := deserializer.ConfirmationIdWithBookingTime(body)
		observers, err := handler.State.ExtendBooking(confirmationId, extendTime)
		if err == nil {
			notification := serializer.NotifyExtend(confirmationId, extendTime)
			observers.Notify(notification)
		}
		reply := serializer.ReplyStatus(err)
		handler.CallingClient.SendMessage(reply)

	case CancelBooking:
		confirmationId := deserializer.ConfirmationId(body)
		observers, err := handler.State.CancelBooking(confirmationId)
		if err == nil {
			notification := serializer.NotifyCancel(confirmationId)
			observers.Notify(notification)
		}
		reply := serializer.ReplyStatus(err)
		handler.CallingClient.SendMessage(reply)

	}
}
