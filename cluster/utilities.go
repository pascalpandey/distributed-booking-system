package cluster

import (
	"log"
	"net"
	"strings"

	"sc4051-server/client"
)

// Check address for filtering of clusterNodes flag
func areSameAddr(a, b *net.UDPAddr) bool {
	// Treat "::" as "127.0.0.1" for demo purposes
	if (((a.IP.String() == "::" && b.IP.String() == "127.0.0.1") || 
		(b.IP.String() == "::" && a.IP.String() == "127.0.0.1") ||
		a.IP.Equal(b.IP)) && a.Port == b.Port) {
		return true
	}
	return false
}

// Filtering of clusterNodes flag
func ExtractClusterServers(str string, conn *net.UDPConn) []*client.Client {
	arr := strings.Split(str, ",")
	res := []*client.Client{}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	for _, addr := range arr {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Printf("Error resolving address: %s", err)
			return nil
		}
		if !areSameAddr(udpAddr, localAddr) {
			res = append(res, &client.Client{Addr: udpAddr, Conn: conn})
		}
	}
	return res
}

// Check if message is a cluster message
func ExtractClusterMsg(str string) (bool, string) {
	substrs := []string{RequestVote, Vote, AcknowledgeState, State, Heartbeat}
	for _, substr := range substrs {
		if strings.Contains(str, substr) {
			return true, substr
		}
	}
	return false, ""
}