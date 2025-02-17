package main

import (
	"log"

	"github.com/distributed-systems-be/handler"
	"github.com/distributed-systems-be/server"
	"github.com/distributed-systems-be/state"
)

var (
	currentState = state.InitState()
)

func main() {
	conn := server.InitUDPServer()
	if conn == nil {
		return
	}

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %s", err)
			continue
		}

		message := string(buffer[:n])
		log.Printf("Received from %s: %s", clientAddr, message)

		handler.HandleMessage(message)

		response := "Message received"
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			log.Printf("Error sending response: %s", err)
		}
	}
}