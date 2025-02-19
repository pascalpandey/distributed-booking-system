package serializer

import (
	"fmt"

	"github.com/google/uuid"
)

func ReplyQueryAvailability(available bool, err error) string {
	if err != nil {
		return fmt.Sprintf("ERROR,%s", err.Error())
	}
	if available {
		return "AVAILABLE"
	}
	return "ERROR,facility already booked for the given period"
}

func ReplyBook(confirmationId uuid.UUID, err error) string {
	if err != nil {
		return fmt.Sprintf("ERROR,%s", err.Error())
	}
	return fmt.Sprintf("SUCCESS,%s", confirmationId)
}

func ReplyStatus(err error) string {
	if err != nil {
		return fmt.Sprintf("ERROR,%s", err.Error())
	}
	return "SUCCESS"
}
