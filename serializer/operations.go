package serializer

import (
	"fmt"
)

// Serializes a response for availability query
func ReplyQueryAvailability(requestId string, available bool, err error) string {
	if err != nil {
		return fmt.Sprintf("%s,ERROR,%s", requestId, err.Error())
	}
	if available {
		return fmt.Sprintf("%s,SUCCESS", requestId)
	}
	return fmt.Sprintf("%s,ERROR,facility already booked for the given period", requestId)
}

// Serializes a response for booking attempt
func ReplyBook(requestId string, confirmationId string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s,ERROR,%s", requestId, err.Error())
	}
	return fmt.Sprintf("%s,SUCCESS,%s", requestId, confirmationId)
}

// Serializes a response for canceling a booking
func ReplyCancel(requestId string, confirmationId string, alreadyCancelled bool) string {
	if alreadyCancelled {
		return fmt.Sprintf("%s,SUCCESS,booking with confirmationId %s already cancelled", requestId, confirmationId)
	}
	return  fmt.Sprintf("%s,SUCCESS", requestId)
}

// Serializes a generic response of success or error with reason
func ReplyStatus(requestId string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s,ERROR,%s", requestId, err.Error())
	}
	return  fmt.Sprintf("%s,SUCCESS", requestId)
}
