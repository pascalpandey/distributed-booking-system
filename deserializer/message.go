package deserializer

import "strings"

func Message(message string) (string, string, []string) {
	lst := strings.Split(message, ",")
	return lst[0], lst[1], lst[2:]
}
