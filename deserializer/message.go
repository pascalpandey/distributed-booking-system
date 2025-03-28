package deserializer

import "strings"

// Deserializes the initial client message to requestId, opcode, and body
func Message(message string) (string, string, []string) {
	lst := strings.Split(message, ",")
	return lst[0], lst[1], lst[2:]
}
