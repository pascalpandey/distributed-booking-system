package main

import (
	"log"
	"time"
	"net"

	"sc4051-server/client"
)

// Forward message to handler server and forward response back to calling client
func forwardMessage(callingClient *client.Client, routerClient *client.Client, message string) error {
	timeout := time.Duration(*timeout * 1000) * time.Millisecond
	response, err := routerClient.SendMessageAndWaitForResponse(message, timeout)
	if err != nil {
		log.Println("Failed to receive reply before timeout")
		return err
	}
	callingClient.SendMessage(response)
	return nil
}

// Forward monitor message and run a background listener for monitor notifications for duration + 3 seconds to account for
// any network latency, listener automatically stops after duration is elapsed
func forwardMonitorMessage(callingClient *client.Client, routerClient *client.Client, message string, duration time.Duration) {
	err := forwardMessage(callingClient, routerClient, message)
	if err != nil {
		return
	}

	listenUntil := time.Now().Add(duration).Add(3 * time.Second)
	go func() {
		log.Printf("Monitor for %s opened \n\n", callingClient.Addr)
		for {
			response, err := routerClient.ListenForResponse(listenUntil)
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			log.Printf("Forwarding monitor notification to %s \n\n", callingClient.Addr)
			callingClient.SendMessage(response)
		}
		log.Printf("Monitor for %s closed \n\n", callingClient.Addr)
	}()
}
