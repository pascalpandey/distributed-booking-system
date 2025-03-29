package main

import (
	"flag"
	"fmt"
	"log"

	"sc4051-server/client"
	"sc4051-server/handler"
	"sc4051-server/serializer"
	"sc4051-server/server"
	"sc4051-server/state"
)

var (
	useCache = flag.Bool("cache", false, "Enable caching")
	useDrop = flag.Bool("drop", false, "Drop every other packet")
	port = flag.Int("port", 9000, "UDP server port")
)

func main() {
	// Parse flags and init state, cache, and handler
	flag.Parse()

	conn := server.InitUDPServer(*port)
	if conn == nil {
		return
	}
	log.Printf("Handler UDP server listening on %d \n\n", *port)

	state := state.InitState()

	var cache map[string]string
	if *useCache {
		cache = map[string]string{}
	}

	handler := handler.Handler{
		State: &state,
		Cache: cache,
	}

	// Used in useDrop=true to indicate whether a packet should be dropped or not
	drop := true

	for {
		// Get raw bytes message from client and convert to string
		buffer := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %s", err)
			continue
		}
		message := string(buffer[:n])

		// If useDrop=true, drop every other packet
		if *useDrop {
			if drop {
				drop = false
				log.Printf("Dropped %s from %s \n\n", message, clientAddr)
				continue
			} else {
				drop = true
			}
		}
		log.Printf("Received from %s: %s", clientAddr, message)
		
		// Handle message based on calling client address
		handler.CallingClient = &client.Client{Conn: conn, Addr: clientAddr}
		cached := handler.HandleMessage(message)
		if cached {
			log.Printf("Message has been cached, resent previous reply \n\n")
			continue
		}
		
		log.Printf("State now:")
        for key, facilityState := range state {
			if len(facilityState.Bookings) > 0 {
				log.Printf("📅 %s BOOKINGS:", key)
				for i, booking := range facilityState.Bookings {
					log.Printf("  [%d] %s/%02d:%02d - %s/%02d:%02d | ID: %s",
						i+1,
						serializer.DayToString(booking.StartTime.Day), 
						booking.StartTime.Hour, 
						booking.StartTime.Minute,
						serializer.DayToString(booking.EndTime.Day), 
						booking.EndTime.Hour, 
						booking.EndTime.Minute,
						booking.ConfirmationId,
					)
				}
			}
			
			if len(facilityState.Observers) > 0 {
				log.Printf("👁 %s OBSERVERS:", key)
				i := 1
				for uuid, observer := range facilityState.Observers {
					log.Printf("  [%d] UUID: %s | Address: %+v", 
						i, 
						uuid, 
						observer.Addr,
					)
					i++
				}
			}
		}
		fmt.Printf("\n")
	}
}
