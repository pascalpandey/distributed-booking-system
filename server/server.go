package server

import (
	"fmt"
	"log"
	"net"
)

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

	log.Printf("UDP server listening on %s", addr)

	return conn
}