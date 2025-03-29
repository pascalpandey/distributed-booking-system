package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"sc4051-server/client"
	"sc4051-server/deserializer"
	"sc4051-server/handler"
	"sc4051-server/server"
)

var (
	port        = flag.Int("port", 9000, "UDP router port")
	trAddr      = flag.String("trAddr", "127.0.0.1:8080", "IP address and port of TR server")
	labAddr     = flag.String("labAddr", "127.0.0.1:8081", "IP address and port of LAB server")
	theatreAddr = flag.String("theatreAddr", "127.0.0.1:8082", "IP address and port of THEATRE server")
	timeout     = flag.Float64("timeout", 1.0, "Router timeout waiting for handler server response")
)

func main() {
	flag.Parse()

	conn := server.InitUDPServer(*port)
	if conn == nil {
		return
	}
	log.Printf("Router UDP server listening on %d \n\n", *port)

	for {
		// Get raw bytes message from client and convert to string
		buffer := make([]byte, 1024)

		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %s", err)
			continue
		}

		message := string(buffer[:n])
		log.Printf("Received from %s: %s", clientAddr, message)

		_, operation, body := deserializer.Message(message)

		// Extract facility type from facility name directly or confirmationId
		// depending on operation type and duration for monitor operation
		var facilityType FacilityType
		var monitorDuration *time.Duration
		switch operation {
		case handler.QueryAvailability, handler.Book:
			facility, _, _ := deserializer.FacilityWithBooking(body)
			facilityType = extractFacilityType(facility)

		case handler.MonitorAvailability:
			facility, duration := deserializer.FacilityWithMonitorDuration(body)
			facilityType = extractFacilityType(facility)
			monitorDuration = &duration

		case handler.ExtendBooking, handler.OffsetBooking:
			confirmationId, _ := deserializer.ConfirmationIdWithBookingTime(body)
			facilityType = extractFacilityType(confirmationId)

		case handler.CancelBooking:
			confirmationId := deserializer.ConfirmationId(body)
			facilityType = extractFacilityType(confirmationId)

		}

		// Create a client to forward message to respective handler server
		var routerConn *net.UDPConn
		var handlerServerAddr *net.UDPAddr
		switch facilityType {
		case TR:
			routerConn, handlerServerAddr = server.InitUDPClient(*trAddr)
		case Lab:
			routerConn, handlerServerAddr = server.InitUDPClient(*labAddr)
		case Theatre:
			routerConn, handlerServerAddr = server.InitUDPClient(*theatreAddr)
		}
		routerClient := &client.Client{Conn: routerConn, Addr: handlerServerAddr}
		log.Printf("Forwarding to: %s", handlerServerAddr)

		// Handle forwarding of monitor and non monitor operations
		callingClient := &client.Client{Conn: conn, Addr: clientAddr}
		if monitorDuration != nil {
			forwardMonitorMessage(callingClient, routerClient, message, *monitorDuration)
		} else {
			forwardMessage(callingClient, routerClient, message)
		}

		fmt.Printf("\n")
	}
}
