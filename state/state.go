package state

import (
	"strconv"
)

type Room = string

type state map[Room][]Booking

func InitState() state {
	state := make(state)
	for i := 1; i <= 10; i++ {
		roomName := "TR" + strconv.Itoa(i)
		state[roomName] = []Booking{}
	}
	return state
}
