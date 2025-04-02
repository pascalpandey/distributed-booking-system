package cluster

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

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

func (clusterState *ClusterState) revertState(backup string) {
	_, dataState := clusterState.deserializeState(backup)
	clusterState.DataState = dataState
}

func (clusterState *ClusterState) SendState(backup string) bool {
	clusterState.ReplicationLock.Lock()
	clusterState.DataState.Id += 1
	clusterState.StateAckCount = 1
	for _, client := range clusterState.ClusterClients {
		client.SendMessage(clusterState.SerializeState())
	}

	time.Sleep(StateInterval)

	majority := int32((len(clusterState.ClusterClients)+1)/2 + 1)
	if atomic.LoadInt32(&clusterState.StateAckCount) < majority {
		log.Printf("Only received %d acks, not enough for commit stage", clusterState.StateAckCount)
		clusterState.ReplicationLock.Unlock()
		return false
	}

	log.Printf("Replication acknowledged with %d acks, starting commit...", clusterState.StateAckCount)

	clusterState.StateAckCount = 1
	for _, client := range clusterState.ClusterClients {
		client.SendMessage(clusterState.SerializeState())
	}

	time.Sleep(StateInterval)

	if atomic.LoadInt32(&clusterState.StateAckCount) < majority {
		log.Printf("Only received %d acks, not enough for commit stage", clusterState.StateAckCount)
		clusterState.revertState(backup)
		clusterState.ReplicationLock.Unlock()
		return false
	}

	clusterState.saveStateToDisk()
	log.Println("Committed")

	clusterState.ReplicationLock.Unlock()
	return true
}

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