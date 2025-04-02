package cluster

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"sc4051-server/client"
	"sc4051-server/serializer"
	"sc4051-server/state"
)

type ClusterState struct {
	ClusterClients   []*client.Client
	DataState        DataState
	PendingDataState DataState
	ServerState      ServerState
	MessageQueue     chan ClusterMessage
	TimeoutReached   bool
	CurrentTerm      int32
	HasVoted         bool
	VoteCount        int32
	LeaderAddr       string
	StateAckCount    int32
	ReplicationLock  sync.Mutex
}

type DataState struct {
	Id    int32
	State state.State
}

type ServerState = int

const (
	Follower ServerState = iota // Represents state as enum of ints starting from 0
	Candidate
	Leader
)

type ClusterMessage struct {
	CallingClient *client.Client
	Opcode        string
	Message       string
}

const (
	Vote              = "VOTE"
	RequestVote       = "REQVOTE"
	State             = "STATE"
	AcknowledgeState  = "ACKSTATE"
	Heartbeat         = "HEARTBEAT"
	ElectionDuration  = 7000 * time.Millisecond
	HeartbeatInterval = 500 * time.Millisecond
	StateInterval     = 200 * time.Millisecond
)

func (clusterState *ClusterState) Start() {
	clusterState.loadState()
	go clusterState.handleMessage()
	for {
		time.Sleep(time.Duration(rand.Intn(2000)+2000) * time.Millisecond)
		if clusterState.TimeoutReached && clusterState.ServerState == Follower {
			clusterState.startElection()
		}
		clusterState.TimeoutReached = true
	}
}

func (clusterState *ClusterState) SerializeState() string {
	encodedData, err := json.Marshal(clusterState.DataState)
	if err != nil {
		log.Printf("Error encoding DataState: %v", err)
		return ""
	}

	data := string(encodedData)

	serialized := fmt.Sprintf("STATE|%d|%s", clusterState.CurrentTerm, data)
	return serialized
}

func (clusterState *ClusterState) deserializeState(str string) (int, DataState) {
	arr := strings.Split(str, "|")
	term, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Printf("failed to parse term: %v", err)
	}

	decodedBytes := []byte(arr[2])
	var decoded DataState
	err = json.Unmarshal(decodedBytes, &decoded)
	if err != nil {
		log.Printf("failed to decode DataState: %v", err)
	}

	return term, decoded
}

func (clusterState *ClusterState) loadState() {
	addr := clusterState.ClusterClients[0].Conn.LocalAddr().(*net.UDPAddr)
	filename := fmt.Sprintf("states/server_%d", addr.Port)
	serialized, err := os.ReadFile(filename)
	if err != nil {
		log.Println("No committed state yet")
		return
	}
	serialized_str := string(serialized)
	term, dataState := clusterState.deserializeState(serialized_str)
	clusterState.CurrentTerm = int32(term)
	clusterState.DataState = dataState
	log.Printf("Loaded state with term %d and log id %d", term, dataState.Id)
}

func (clusterState *ClusterState) saveStateToDisk() {
	serialized := clusterState.SerializeState()
	addr := clusterState.ClusterClients[0].Conn.LocalAddr().(*net.UDPAddr)
	filename := fmt.Sprintf("states/server_%d", addr.Port)
	err := os.WriteFile(filename, []byte(serialized), 0644)
	if err != nil {
		log.Printf("Failed to save state to %s: %+v", filename, err)
	}
}

func (clusterState *ClusterState) logDataState() {
	log.Printf("State now with log id %d:", clusterState.DataState.Id)
	for key, facilityState := range clusterState.DataState.State {
		if len(facilityState.Bookings) > 0 {
			log.Printf("📅 %s BOOKINGS:", key)
			for i, booking := range facilityState.Bookings {
				log.Printf("  [%d] %s/%02d:%02d - %s/%02d:%02d | ID: %s",
					i+1,
					serializer.DayToString(booking.StartTime.Day),
					booking.StartTime.Hour,
					booking.StartTime.Minute,
					serializer.DayToString(booking.EndTime.Day),
					booking.EndTime.Hour,
					booking.EndTime.Minute,
					booking.ConfirmationId,
				)
			}
		}

		if len(facilityState.Observers) > 0 {
			log.Printf("👁 %s OBSERVERS:", key)
			i := 1
			for uuid, observer := range facilityState.Observers {
				log.Printf("  [%d] UUID: %s | Address: %+v",
					i,
					uuid,
					observer.Addr,
				)
				i++
			}
		}
	}
	fmt.Printf("\n")
}
