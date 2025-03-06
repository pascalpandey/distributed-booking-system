package main

import (
	"flag"
	"log"

	"sc4051-server/client"
	"sc4051-server/handler"
	"sc4051-server/serializer"
	"sc4051-server/server"
	"sc4051-server/state"
)

func main() {
	useCache := flag.Bool("cache", false, "Enable caching")
	useDrop := flag.Bool("drop", false, "Drop every other packet")
	port := flag.Int("port", 8080, "UDP server port")
	flag.Parse()

	conn := server.InitUDPServer(*port)
	if conn == nil {
		return
	}

	state := state.InitState()

	var cache map[string]string
	if *useCache {
		cache = map[string]string{}
	}

	for {
		buffer := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %s", err)
			continue
		}

		message := string(buffer[:n])
		drop := true
		if *useDrop {
			if drop {
				drop = false
				log.Printf("Dropped %s from %s", clientAddr, message)
				continue
			} else {
				drop = true
			}
		}
		log.Printf("Received from %s: %s", clientAddr, message)

		handler := handler.Handler{
			State: &state,
			CallingClient: &client.Client{
				Conn: conn,
				Addr: clientAddr,
			},
			Cache: cache,
		}
		handler.HandleMessage(message)
		
		log.Println("State now:")
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
	}
}
