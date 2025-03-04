package main

import (
	"flag"
	"log"

	"github.com/distributed-systems-be/client"
	"github.com/distributed-systems-be/handler"
	"github.com/distributed-systems-be/serializer"
	"github.com/distributed-systems-be/server"
	"github.com/distributed-systems-be/state"
)

func main() {
	useCache := flag.Bool("cache", false, "Enable caching")
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

        log.Println("State Now:")
        for key, facilityState := range state {
            log.Printf("%s bookings: ", key)
            for _, booking := range facilityState.Bookings {
                log.Printf("[%s/%02d/%02d - %s/%02d/%02d, %s] ", 
                    serializer.DayToString(booking.StartTime.Day), booking.StartTime.Hour, booking.StartTime.Minute,
                    serializer.DayToString(booking.EndTime.Day), booking.EndTime.Hour, booking.EndTime.Minute,
                    booking.ConfirmationId.String(),
                )
            }
            log.Printf("%s observers: ", key)
            for uuid, observer := range facilityState.Observers {
                log.Printf("(%s -> %+v) ", uuid, observer)
            }
        }
	}
}
