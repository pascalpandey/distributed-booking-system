package handler

import (
	"fmt"
	"log"
	"sc4051-server/client"
	"sc4051-server/cluster"
	"sc4051-server/deserializer"
	"sc4051-server/serializer"
	"sc4051-server/state"
)

const (
	QueryAvailability   = "QUERY"   // Checking availability of a facility
	Book                = "BOOK"    // Booking a facility
	OffsetBooking       = "OFFSET"  // Offsetting the start time of a booking
	MonitorAvailability = "MONITOR" // Register observer to monitor facility state
	ExtendBooking       = "EXTEND"  // Extending the end time of a booking
	CancelBooking       = "CANCEL"  // Cancelling a booking
)

type Handler struct {
	State         *state.State          // Reference to the current state of the server
	CallingClient *client.Client        // The client that initiated the request
	Cache         map[string]string     // Cache for storing responses to avoid redundant processing
	ClusterState  *cluster.ClusterState // ClusterState when using cluster setup
}

// Handles incoming messages based on the type of operation in the message.
// First the initial message is deserialized to its requestId, operation, and main body.
// If useCache=true, requestId is checked and if it has been cached, immediately return true and resend the previous response.
// Otherwise, the handler follows a general process of deserializing the main body, handling the operation, for certain
// operations, notify observers monitoring the facility, serialize the response, and if useCache=true save the response
// to the cache before sending it back to the calling client.
func (handler *Handler) HandleMessage(message string, dropReply bool) {
	requestId, operation, body := deserializer.Message(message)

	if handler.Cache != nil {
		reply, requestHandled := handler.Cache[requestId]
		if requestHandled {
			handler.CallingClient.SendMessage(reply)
			log.Printf("Message has been cached, resent previous reply \n\n")
			return
		}
	}

	var backup string
	if handler.ClusterState != nil {
		backup = handler.ClusterState.SerializeState()
	}

	switch operation {
	case QueryAvailability:
		facility, startTime, endTime := deserializer.FacilityWithBooking(body)
		available, err := handler.State.QueryAvailability(facility, startTime, endTime)
		reply := serializer.ReplyQueryAvailability(requestId, available, err)
		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		if !dropReply {
			handler.CallingClient.SendMessage(reply)
		}

	case Book:
		facility, startTime, endTime := deserializer.FacilityWithBooking(body)
		observers, confirmationId, err := handler.State.Book(facility, startTime, endTime)

		if err == nil {
			notification := serializer.NotifyBook(facility, requestId, confirmationId, startTime, endTime)
			observers.Notify(notification)
		}

		if handler.ClusterState != nil {
			success := handler.ClusterState.SendState(backup)
			if !success {
				log.Printf("Replication failed, reset state")
				handler.LogState()
				return
			}
		}

		reply := serializer.ReplyBook(requestId, confirmationId, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		if !dropReply {
			handler.CallingClient.SendMessage(reply)
		}

	case OffsetBooking:
		confirmationId, offsetTime := deserializer.ConfirmationIdWithBookingTime(body)
		observers, err := handler.State.OffsetBooking(confirmationId, offsetTime)

		if handler.ClusterState != nil {
			success := handler.ClusterState.SendState(backup)
			if !success {
				log.Printf("Replication failed, reset state")
				handler.LogState()
				return
			}
		}

		if err == nil {
			notification := serializer.NotifyOffset(confirmationId, offsetTime)
			observers.Notify(notification)
		}

		reply := serializer.ReplyStatus(requestId, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		if !dropReply {
			handler.CallingClient.SendMessage(reply)
		}

	case MonitorAvailability:
		facility, monitorDuration := deserializer.FacilityWithMonitorDuration(body)
		err := handler.State.Monitor(handler.CallingClient, facility, monitorDuration)
		reply := serializer.ReplyStatus(requestId, err)

		if handler.ClusterState != nil {
			success := handler.ClusterState.SendState(backup)
			if !success {
				log.Printf("Replication failed, reset state")
				handler.LogState()
				return
			}
		}

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		if !dropReply {
			handler.CallingClient.SendMessage(reply)
		}

	case ExtendBooking:
		confirmationId, extendTime := deserializer.ConfirmationIdWithBookingTime(body)
		observers, err := handler.State.ExtendBooking(confirmationId, extendTime)

		if handler.ClusterState != nil {
			success := handler.ClusterState.SendState(backup)
			if !success {
				log.Printf("Replication failed, reset state")
				handler.LogState()
				return
			}
		}

		if err == nil {
			notification := serializer.NotifyExtend(confirmationId, extendTime)
			observers.Notify(notification)
		}

		reply := serializer.ReplyStatus(requestId, err)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		if !dropReply {
			handler.CallingClient.SendMessage(reply)
		}

	case CancelBooking:
		confirmationId := deserializer.ConfirmationId(body)
		observers, alreadyCancelled := handler.State.CancelBooking(confirmationId)

		if handler.ClusterState != nil {
			success := handler.ClusterState.SendState(backup)
			if !success {
				log.Printf("Replication failed, reset state")
				handler.LogState()
				return
			}
		}

		if !alreadyCancelled {
			notification := serializer.NotifyCancel(confirmationId)
			observers.Notify(notification)
		}

		reply := serializer.ReplyCancel(requestId, confirmationId, alreadyCancelled)

		if handler.Cache != nil {
			handler.Cache[requestId] = reply
		}

		if !dropReply {
			handler.CallingClient.SendMessage(reply)
		}
	}

	handler.LogState()
}

func (handler *Handler) LogState() {
	log.Printf("State now:")
	for key, facilityState := range *handler.State {
		if len(facilityState.Bookings) > 0 {
			log.Printf("📅 %s BOOKINGS:", key)
			for i, booking := range facilityState.Bookings {
				log.Printf("  [%d] %s/%02d:%02d - %s/%02d:%02d | ID: %s",
					i+1,
					serializer.DayToString(booking.StartTime.Day),
					booking.StartTime.Hour,
					booking.StartTime.Minute,
					serializer.DayToString(booking.EndTime.Day),
					booking.EndTime.Hour,
					booking.EndTime.Minute,
					booking.ConfirmationId,
				)
			}
		}

		if len(facilityState.Observers) > 0 {
			log.Printf("👁 %s OBSERVERS:", key)
			i := 1
			for uuid, observer := range facilityState.Observers {
				log.Printf("  [%d] UUID: %s | Address: %+v",
					i,
					uuid,
					observer.Addr,
				)
				i++
			}
		}
	}
	fmt.Printf("\n")
}
