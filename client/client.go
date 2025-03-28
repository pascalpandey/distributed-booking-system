package client

import (
	"net"
)

type Client struct {
	Conn *net.UDPConn // UDP connection used by the server to send messages back
	Addr *net.UDPAddr // UDP address of client
}

// Sends a message to the client's address via UDP
func (client *Client) SendMessage(message string) {
	client.Conn.WriteToUDP([]byte(message), client.Addr)
}
