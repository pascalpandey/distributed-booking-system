package client

import (
	"net"
)

type Client struct {
	Conn *net.UDPConn
	Addr *net.UDPAddr
}

func (client *Client) SendMessage(message string) {
	client.Conn.WriteToUDP([]byte(message), client.Addr)
}
