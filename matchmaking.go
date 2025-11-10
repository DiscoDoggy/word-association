package main

import "errors"

type PlayerQueue struct {
	queue []*Client
}

func CreatePlayerQueue() *PlayerQueue {
	return &PlayerQueue{
		queue: make([]*Client, 0, 2),
	}	
}

func (pq *PlayerQueue) addClientToQueue(client *Client) error {
	for _, val := range pq.queue {
		if client == val {
			return errors.New("Client is already queued")
		}
	}

	pq.queue = append(pq.queue, client)
	return nil
}

func (pq *PlayerQueue) removePlayerFromQueue(client *Client) error{
	for idx, c := range pq.queue {
		if c == client {
			pq.queue = append(pq.queue[:idx], pq.queue[idx+1:]...)	
			break
		} 
	}

	return nil
}

func (pq *PlayerQueue) ScanPlayerQueue() {
	defer func() {
		pq.queue = pq.queue[:0]
	}()

	for {
		// 2 here represents the two players per match
		if len(pq.queue) % 2 == 0 {
			matchPlayers := pq.queue[:2]
			//does this need to go and start its own go routine???
			CreateMatch(matchPlayers)
			if len(pq.queue) == 2 {
				pq.queue = pq.queue[:0]
			} else {
				pq.queue = pq.queue[2:] 
			}
		}
	}
}

