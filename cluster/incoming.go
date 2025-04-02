package cluster

import (
	"log"
	"strconv"
	"strings"
	"sync/atomic"
)

// Handle state updates in replication process
func (clusterState *ClusterState) handleState(message ClusterMessage) {
	switch clusterState.ServerState {
	case Leader:
		return
	case Candidate:
		return
	case Follower:
		term, dataState := clusterState.deserializeState(message.Message)
		if dataState.Id < clusterState.DataState.Id {
			log.Printf("Outdated state, ignored")
			return
		}
		if clusterState.CurrentTerm < int32(term) {
			clusterState.CurrentTerm = int32(term)
		}

		log.Printf("Received state with term %d and id %d", term, dataState.Id)
		log.Printf("Current pending state has term %d and id %d", clusterState.CurrentTerm, clusterState.PendingDataState.Id)

		// pending data state is sent again, then commit to disk and make it current data state
		if clusterState.PendingDataState.Id == dataState.Id && clusterState.CurrentTerm == int32(term) {
			log.Printf("Detected resent pending state, commiting state with term %d and id %d...", clusterState.CurrentTerm, dataState.Id)
			clusterState.DataState = dataState
			clusterState.saveStateToDisk()
			clusterState.logDataState()
		// else keep it pending until leader sends another message
		} else {
			log.Printf("Detected new state, saving to state with term %d and id %d...", clusterState.CurrentTerm, dataState.Id)
			clusterState.PendingDataState = dataState
		}

		message.CallingClient.SendMessage("ACKSTATE")
	}
}

// Increment counter of state acknowledgements
func (clusterState *ClusterState) handleAckState() {
	switch clusterState.ServerState {
	case Leader:
		atomic.AddInt32(&clusterState.StateAckCount, 1)
	case Candidate:
		return
	case Follower:
		return
	}
}

// Update timeout after heartbeat
func (clusterState *ClusterState) handleHeartbeat(message ClusterMessage) {
	clusterState.LeaderAddr = message.CallingClient.Addr.String()
	clusterState.TimeoutReached = false
}

// Perform checks if a vote request should be given a vote
func (clusterState *ClusterState) handleReqVote(message ClusterMessage) {
	arr := strings.Split(message.Message, ",")
	term, _ := strconv.Atoi(arr[1])
	logId, _ := strconv.Atoi(arr[2])
	if clusterState.DataState.Id > int32(logId) {
		if clusterState.CurrentTerm < int32(term) {
			clusterState.CurrentTerm = int32(term)
		}
		log.Println("Ignored REQVOTE as log id is out of date")
		return
	}
	switch clusterState.ServerState {
	case Candidate, Leader:
		if clusterState.CurrentTerm < int32(term) {
			clusterState.ServerState = Follower
			clusterState.CurrentTerm = int32(term)
			clusterState.HasVoted = true
			message.CallingClient.SendMessage("VOTE")
		}
	case Follower:
		if !clusterState.HasVoted || clusterState.CurrentTerm < int32(term) {
			clusterState.CurrentTerm = int32(term)
			clusterState.HasVoted = true
			message.CallingClient.SendMessage("VOTE")
		}
	}
}

// Update vote count if receive vote
func (clusterState *ClusterState) handleVote() {
	switch clusterState.ServerState {
	case Leader:
		return
	case Candidate:
		atomic.AddInt32(&clusterState.VoteCount, 1)
		majority := int32((len(clusterState.ClusterClients)+1)/2 + 1)
		if atomic.LoadInt32(&clusterState.VoteCount) >= majority {
			log.Printf("Elected with %d votes", clusterState.VoteCount)
			clusterState.ServerState = Leader
			go clusterState.sendHeartbeats()
		}
	case Follower:
		return
	}
}

// Main loop to handle cluster requests
func (clusterState *ClusterState) handleMessage() {
	for message := range clusterState.MessageQueue {
		switch message.Opcode {
		case Vote:
			clusterState.handleVote()
		case RequestVote:
			clusterState.handleReqVote(message)
		case State:
			clusterState.handleState(message)
		case AcknowledgeState:
			clusterState.handleAckState()
		case Heartbeat:
			clusterState.handleHeartbeat(message)
		}
	}
}
