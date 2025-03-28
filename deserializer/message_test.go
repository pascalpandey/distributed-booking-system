package deserializer

import (
	"testing"
)

func TestMessage(t *testing.T) {
	tests := []struct {
		message           string
		expectedRequestId string
		expectedOpcode    string
		expectedBody      []string
	}{
		{
			message:           "123,BOOK,venue1,venue2",
			expectedRequestId: "123",
			expectedOpcode:    "BOOK",
			expectedBody:      []string{"venue1", "venue2"},
		},
		{
			message:           "456,CANCEL",
			expectedRequestId: "456",
			expectedOpcode:    "CANCEL",
			expectedBody:      []string{},
		},
	}

	for _, test := range tests {
		reqId, opcode, body := Message(test.message)

		if reqId != test.expectedRequestId || opcode != test.expectedOpcode || !equal(body, test.expectedBody) {
			t.Errorf("For message %v, expected (%v, %v, %v), but got (%v, %v, %v)",
				test.message, test.expectedRequestId, test.expectedOpcode, test.expectedBody, reqId, opcode, body)
		}
	}
}

// Helper function to compare slices
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
