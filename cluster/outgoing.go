package cluster

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

// Send heartbeats periodically if a leader
func (clusterState *ClusterState) sendHeartbeats() {
	for {
		if clusterState.ServerState != Leader {
			return
		}
		for _, client := range clusterState.ClusterClients {
			client.SendMessage("HEARTBEAT")
		}
		time.Sleep(HeartbeatInterval)
	}
}

// Revert state incase of send state failure
func (clusterState *ClusterState) revertState(backup string) {
	_, dataState := clusterState.deserializeState(backup)
	clusterState.DataState = dataState
}

// Send state to followers for replication
func (clusterState *ClusterState) SendState(backup string) bool {
	clusterState.ReplicationLock.Lock()

	log.Println("Starting replication, broadcasting state...")

	clusterState.DataState.Id += 1
	clusterState.StateAckCount = 1
	for _, client := range clusterState.ClusterClients {
		client.SendMessage(clusterState.SerializeState())
	}

	time.Sleep(StateInterval)

	majority := int32((len(clusterState.ClusterClients)+1)/2 + 1)
	if atomic.LoadInt32(&clusterState.StateAckCount) < majority {
		log.Printf("Only received %d acks, not enough for commit stage", clusterState.StateAckCount)
		clusterState.revertState(backup)
		clusterState.ReplicationLock.Unlock()
		return false
	}

	log.Printf("Replication acknowledged with %d acks, starting commit...", clusterState.StateAckCount)

	clusterState.saveStateToDisk()

	clusterState.StateAckCount = 1
	for _, client := range clusterState.ClusterClients {
		client.SendMessage(clusterState.SerializeState())
	}

	log.Println("Committed")

	clusterState.ReplicationLock.Unlock()
	return true
}

// Start election process after timeout
func (clusterState *ClusterState) startElection() {
	clusterState.CurrentTerm += 1
	log.Printf("Starting election with term %d", clusterState.CurrentTerm)

	clusterState.VoteCount = 1
	clusterState.HasVoted = true
	clusterState.ServerState = Candidate
	clusterState.LeaderAddr = ""
	for _, client := range clusterState.ClusterClients {
		client.SendMessage(fmt.Sprintf("REQVOTE,%d,%d", clusterState.CurrentTerm, clusterState.DataState.Id))
	}

	time.Sleep(ElectionDuration)

	if clusterState.ServerState != Leader {
		clusterState.ServerState = Follower
		log.Printf("Only got %d votes, election ended", clusterState.VoteCount)
	}
}