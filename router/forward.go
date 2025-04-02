package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"sc4051-server/client"
)

// Check if redirection required for cluster setup
func checkRedirect(message string) (bool ,string) {
	if strings.Contains(message, "REDIRECT") {
		arr := strings.Split(message, ",")
		return true, arr[1]
	}
	return false, ""
}

// Forward message to handler server and forward response back to calling client, check redirection in case using cluster setup
func forwardMessage(callingClient *client.Client, routerClient *client.Client, addrs []string, message string) error {
	timeout := time.Duration(*timeout * 1000) * time.Millisecond
	for _, addr := range(addrs) {
		testAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Printf("Failed to resolve new address: %+v", err)
			return err
		}

		routerClient.Addr = testAddr
		response, err := routerClient.SendMessageAndWaitForResponse(message, timeout)
		if err != nil {
			log.Println("Failed to receive reply before timeout, retrying other servers")
			continue
		}

		if isRedirect, newAddr := checkRedirect(response); isRedirect {
			destAddr, err := net.ResolveUDPAddr("udp", newAddr)
			if err != nil {
				log.Printf("Failed to resolve new address: %+v", err)
				return err
			}
			routerClient.Addr = destAddr
			log.Printf("Redirected to: %s", routerClient.Addr)
			response, err := routerClient.SendMessageAndWaitForResponse(message, timeout)
			if err != nil {
				log.Println("Failed to receive reply before timeout, retrying other servers")
				continue
			}
			callingClient.SendMessage(response)
			log.Printf("Sent reply %s to %s", response, callingClient.Addr)
			return nil
		} else {
			callingClient.SendMessage(response)
			log.Printf("Sent reply %s to %s", response, callingClient.Addr)
			return nil
		}
	}
	return fmt.Errorf("failed to reach any servers")
}

// Forward monitor message and run a background listener for monitor notifications for duration + 3 seconds to account for
// any network latency, listener automatically stops after duration is elapsed
func forwardMonitorMessage(callingClient *client.Client, routerClient *client.Client, addrs []string, message string, duration time.Duration) {
	err := forwardMessage(callingClient, routerClient, addrs, message)
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
