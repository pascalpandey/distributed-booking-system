package serializer

import (
	"fmt"
	"testing"
)

func TestReplyQueryAvailability(t *testing.T) {
	tests := []struct {
		requestId string
		available bool
		err       error
		expected  string
	}{
		{
			requestId: "12345",
			available: true,
			err:       nil,
			expected:  "12345,SUCCESS",
		},
		{
			requestId: "12345",
			available: false,
			err:       fmt.Errorf("facility already booked for the given period"),
			expected:  "12345,ERROR,facility already booked for the given period",
		},
	}

	for _, test := range tests {
		result := ReplyQueryAvailability(test.requestId, test.available, test.err)
		if result != test.expected {
			t.Errorf("For params %s, %v, %+v, expected %v, but got %v", test.requestId, test.available, test.err, test.expected, result)
		}
	}
}

func TestReplyBook(t *testing.T) {
	tests := []struct {
		requestId      string
		confirmationId string
		err            error
		expected       string
	}{
		{
			requestId:      "12345",
			confirmationId: "67890",
			err:            nil,
			expected:       "12345,SUCCESS,67890",
		},
		{
			requestId:      "12345",
			confirmationId: "67890",
			err:            fmt.Errorf("booking failed"),
			expected:       "12345,ERROR,booking failed",
		},
	}

	for _, test := range tests {
		result := ReplyBook(test.requestId, test.confirmationId, test.err)
		if result != test.expected {
			t.Errorf("For params %s, %s, %+v, expected %+v, but got %v", test.requestId, test.confirmationId, test.err, test.expected, result)
		}
	}
}

func TestReplyCancel(t *testing.T) {
	tests := []struct {
		requestId        string
		confirmationId   string
		alreadyCancelled bool
		expected         string
	}{
		{
			requestId:        "12345",
			confirmationId:   "67890",
			alreadyCancelled: false,
			expected:         "12345,SUCCESS",
		},
		{
			requestId:        "12345",
			confirmationId:   "67890",
			alreadyCancelled: true,
			expected:         "12345,SUCCESS,booking with confirmationId 67890 already cancelled",
		},
	}

	for _, test := range tests {
		result := ReplyCancel(test.requestId, test.confirmationId, test.alreadyCancelled)
		if result != test.expected {
			t.Errorf("For params %s, %s, %v, expected %v, but got %v", test.requestId, test.confirmationId, test.alreadyCancelled, test.expected, result)
		}
	}
}

func TestReplyStatus(t *testing.T) {
	tests := []struct {
		requestId string
		err       error
		expected  string
	}{
		{
			requestId: "12345",
			err:       nil,
			expected:  "12345,SUCCESS",
		},
		{
			requestId: "12345",
			err:       fmt.Errorf("some error occurred"),
			expected:  "12345,ERROR,some error occurred",
		},
	}

	for _, test := range tests {
		result := ReplyStatus(test.requestId, test.err)
		if result != test.expected {
			t.Errorf("For params %s, %+v, expected %v, but got %v", test.requestId, test.err, test.expected, result)
		}
	}
}
