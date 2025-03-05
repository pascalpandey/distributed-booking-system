package serializer

import (
	"fmt"
)

func ReplyQueryAvailability(requestId string, available bool, err error) string {
	if err != nil {
		return fmt.Sprintf("ERROR,%s", err.Error())
	}
	if available {
		return fmt.Sprintf("%s,AVAILABLE", requestId)
	}
	return fmt.Sprintf("%s,ERROR,facility already booked for the given period", requestId)
}

func ReplyBook(requestId string, confirmationId string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s,ERROR,%s", requestId, err.Error())
	}
	return fmt.Sprintf("%s,SUCCESS,%s", requestId, confirmationId)
}

func ReplyCancel(requestId string, confirmationId string, alreadyCancelled bool) string {
	if alreadyCancelled {
		return fmt.Sprintf("%s,SUCCESS,booking with confirmationId %s already cancelled", requestId, confirmationId)
	}
	return  fmt.Sprintf("%s,SUCCESS", requestId)
}

func ReplyStatus(requestId string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s,ERROR,%s", requestId, err.Error())
	}
	return  fmt.Sprintf("%s,SUCCESS", requestId)
}
