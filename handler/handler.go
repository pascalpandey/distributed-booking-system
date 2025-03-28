package handler

import (
	"sc4051-server/client"
	"sc4051-server/deserializer"
	"sc4051-server/serializer"
	"sc4051-server/state"
)

const (
	QueryAvailability   = "QUERY"    // Checking availability of a facility
	Book                = "BOOK"     // Booking a facility
	OffsetBooking       = "OFFSET"   // Offsetting the start time of a booking
	MonitorAvailability = "MONITOR"  // Register observer to monitor facility state
	ExtendBooking       = "EXTEND"   // Extending the end time of a booking
	CancelBooking       = "CANCEL"   // Cancelling a booking
)

type Handler struct {
	State         *state.State        // Reference to the current state of the server
	CallingClient *client.Client      // The client that initiated the request
	Cache         map[string]string   // Cache for storing responses to avoid redundant processing
}

// Handles incoming messages based on the type of operation in the message.
// First the initial message is deserialized to its requestId, operation, and main body.
// If useCache=true, requestId is checked and if it has been cached, immediately return true and resend the previous response.
// Otherwise, the handler follows a general process of deserializing the main body, handling the operation, for certain
// operations, notify observers monitoring the facility, serialize the response, and if useCache=true save the response
// to the cache before sending it back to the calling client and return false to indicate this message isn't handled by the cache.
func (handler *Handler) HandleMessage(message string) (bool) {
	requestId, operation, body := deserializer.Message(message)

	if handler.Cache != nil {
		reply, requestHandled := handler.Cache[requestId]
		if requestHandled {
			handler.CallingClient.SendMessage(reply)
			return true
		}
	}

	switch operation {
	case QueryAvailability:
		facility, startTime, endTime := deserializer.FacilityWithBooking(body)
		available, err := handler.State.QueryAvailability(facility, startTime, endTime)
		reply := serializer.ReplyQueryAvailability(requestId, available, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		handler.CallingClient.SendMessage(reply)

	case Book:
		facility, startTime, endTime := deserializer.FacilityWithBooking(body)
		observers, confirmationId, err := handler.State.Book(facility, startTime, endTime)

		if err == nil {
			notification := serializer.NotifyBook(facility, requestId, confirmationId, startTime, endTime)
			observers.Notify(notification)
		}

		reply := serializer.ReplyBook(requestId, confirmationId, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		handler.CallingClient.SendMessage(reply)

	case OffsetBooking:
		confirmationId, offsetTime := deserializer.ConfirmationIdWithBookingTime(body)
		observers, err := handler.State.OffsetBooking(confirmationId, offsetTime)

		if err == nil {
			notification := serializer.NotifyOffset(confirmationId, offsetTime)
			observers.Notify(notification)
		}

		reply := serializer.ReplyStatus(requestId, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		handler.CallingClient.SendMessage(reply)

	case MonitorAvailability:
		facility, monitorDuration := deserializer.FacilityWithMonitorDuration(body)
		err := handler.State.Monitor(handler.CallingClient, facility, monitorDuration)
		reply := serializer.ReplyStatus(requestId, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		handler.CallingClient.SendMessage(reply)

	case ExtendBooking:
		confirmationId, extendTime := deserializer.ConfirmationIdWithBookingTime(body)
		observers, err := handler.State.ExtendBooking(confirmationId, extendTime)

		if err == nil {
			notification := serializer.NotifyExtend(confirmationId, extendTime)
			observers.Notify(notification)
		}

		reply := serializer.ReplyStatus(requestId, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		handler.CallingClient.SendMessage(reply)

	case CancelBooking:
		confirmationId := deserializer.ConfirmationId(body)
		observers, alreadyCancelled := handler.State.CancelBooking(confirmationId)

		if !alreadyCancelled {
			notification := serializer.NotifyCancel(confirmationId)
			observers.Notify(notification)
		}

		reply := serializer.ReplyCancel(requestId, confirmationId, alreadyCancelled)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		handler.CallingClient.SendMessage(reply)
	}

	return false
}
