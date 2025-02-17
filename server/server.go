package server

import (
	"net"
	"log"
)

func InitUDPServer() *net.UDPConn {
	addr := ":8080"
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
	defer conn.Close()

	log.Printf("UDP server listening on %s", addr)

	return conn
}