package client

import (
	"log"
	"net"
	"time"
)

type Client struct {
	Conn *net.UDPConn // UDP connection used by the server to send messages back
	Addr *net.UDPAddr // UDP address of destination (calling client if server, handler server if load balancer)
}

// Sends a message to the client's address via UDP
func (client *Client) SendMessage(message string) {
	_, err := client.Conn.WriteToUDP([]byte(message), client.Addr)
	if err != nil {
		log.Printf("Failed to send message: %+v", err)
	}
}

// Sends a message to the destination address using established connection
func (client *Client) SendMessageAndWaitForResponse(message string, timeout time.Duration) (string, error) {
	client.SendMessage(message)

	client.Conn.SetReadDeadline(time.Now().Add(timeout))
	buffer := make([]byte, 1024)
	n, _, err := client.Conn.ReadFromUDP(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}

// Listens for responses from the destination address until designated time
func (client *Client) ListenForResponse(until time.Time) (string, error) {
	client.Conn.SetReadDeadline(until)
	buffer := make([]byte, 1024)
	n, _, err := client.Conn.ReadFromUDP(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}
