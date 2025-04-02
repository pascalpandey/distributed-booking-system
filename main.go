package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"sc4051-server/client"
	"sc4051-server/cluster"
	"sc4051-server/handler"
	"sc4051-server/server"
	"sc4051-server/state"
)

var (
	useCache     = flag.Bool("cache", false, "Enable caching")
	useDrop      = flag.Bool("drop", false, "Drop every other packet")
	port         = flag.Int("port", 9000, "UDP server port")
	clusterNodes = flag.String("cluster", "", "List of ordered IP addresses of servers in the cluster")
)

func main() {
	// Parse flags and init state, cache, and handler
	flag.Parse()

	conn := server.InitUDPServer(*port)
	if conn == nil {
		return
	}
	log.Printf("Handler UDP server listening on %d \n\n", *port)

	state := state.InitState()

	var cache map[string]string
	if *useCache {
		cache = map[string]string{}
	}

	handler := handler.Handler{
		State: &state,
		Cache: cache,
	}

	// Init cluster state if using cluster setup
	var clusterState *cluster.ClusterState
	if *clusterNodes != "" {
		clusterState = &cluster.ClusterState{
			ClusterClients:   cluster.ExtractClusterServers(*clusterNodes, conn),
			DataState:        cluster.DataState{Id: 0, State: state},
			PendingDataState: cluster.DataState{Id: 0, State: state},
			ServerState:      cluster.Follower,
			MessageQueue:     make(chan cluster.ClusterMessage),
			TimeoutReached:   true,
			CurrentTerm:      0,
			HasVoted:         false,
			VoteCount:        0,
			LeaderAddr:       "",
			StateAckCount:    0,
			ReplicationLock:  sync.Mutex{},
		}
		handler.State = &clusterState.DataState.State
		handler.ClusterState = clusterState
		go clusterState.Start()
	}

	// Used in useDrop=true to indicate whether a packet should be dropped or not
	drop := false

	for {
		// Get raw bytes message from client and convert to string
		buffer := make([]byte, 65536)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %s", err)
			continue
		}
		message := string(buffer[:n])
		callingClient := &client.Client{Conn: conn, Addr: clientAddr}

		log.Printf("Received from %s: %s", clientAddr, message)

		// Append cluster messages and skip handling if not leader
		if *clusterNodes != "" {
			if isClusterMsg, opcode := cluster.ExtractClusterMsg(message); isClusterMsg {
				clusterState.MessageQueue <- cluster.ClusterMessage{
					Message:       message,
					CallingClient: callingClient,
					Opcode:        opcode,
				}
				continue
			} else if clusterState.ServerState != cluster.Leader {
				log.Printf("Redirected to %s", clusterState.LeaderAddr)
				callingClient.SendMessage(fmt.Sprintf("REDIRECT,%s", clusterState.LeaderAddr))
				continue
			}
		}

		// If useDrop=true, drop every other packet
		if *useDrop {
			if drop {
				drop = false
				log.Printf("Will drop reply to %s from %s \n\n", message, clientAddr)
			} else {
				drop = true
			}
		}

		// Handle message based on calling client address
		handler.CallingClient = callingClient
		go handler.HandleMessage(message, drop)
	}
}
