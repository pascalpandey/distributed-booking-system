package server

import (
	"fmt"
	"log"
	"net"
)

// Initializes the UDP server on the given port
func InitUDPServer(port int) *net.UDPConn {
	addr := fmt.Sprintf(":%d", port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Printf("Error resolving address: %s", err)
		return nil
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Error starting UDP server: %s", err)
		return nil
	}

	return conn
}

// Initializes a UDP client for to send messages to other servers,
// this opens a socket on a random port on the server
func InitUDPClient(destAddrStr string) (*net.UDPConn, *net.UDPAddr) {
	destAddr, err := net.ResolveUDPAddr("udp", destAddrStr)
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}

	sourceAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}

	conn, err := net.ListenUDP("udp", sourceAddr)
	if err != nil {
		log.Printf("Error starting UDP client: %s", err)
		return nil, nil
	}

	return conn, destAddr
}
