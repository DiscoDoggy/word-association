package main

import (
	"errors"
	"log"
	"sync"
)

type PlayerQueue struct {
	queue []*Client
	addCh chan *Client
	removeCh chan *Client

	manager 	*Manager
	sync.RWMutex
}

func CreatePlayerQueue(manager *Manager) *PlayerQueue {
	return &PlayerQueue{
		queue: make([]*Client, 0, 2),
		addCh: make(chan *Client),
		removeCh: make(chan *Client),
		manager: manager,
	}	
}

func (pq *PlayerQueue) addClientToQueue(client *Client) error {
	pq.Lock()
	defer pq.Unlock()

	for _, val := range pq.queue {
		if client == val {
			return errors.New("Client is already queued")
		}
	}

	pq.queue = append(pq.queue, client)
	return nil
}

func (pq *PlayerQueue) removePlayerFromQueue(client *Client) error{
	pq.Lock()
	defer pq.Unlock()

	for idx, c := range pq.queue {
		if c == client {
			pq.queue = append(pq.queue[:idx], pq.queue[idx+1:]...)	
			break
		} 
	}

	return nil
}

func (pq *PlayerQueue) Len() int {
	pq.Lock()
	defer pq.Unlock()

	return len(pq.queue)
}

func (pq *PlayerQueue) popNewMatchPlayers(numPlayers int) []*Client {
	pq.Lock()
	defer pq.Unlock()

	players := pq.queue[:numPlayers]
	pq.queue = pq.queue[numPlayers:]

	return players
}

func (pq *PlayerQueue) ScanPlayerQueue() {
	defer func() {
		for _, client := range pq.queue {
			err := pq.removePlayerFromQueue(client)
			if err != nil {
				log.Println("issue removing client from queue:", err)
			}
		}

		pq.queue = nil
	}()

	for {
		select {
		case client, ok := <-pq.addCh:
			if !ok {
				log.Println("issue reading client from add client channel:")
				return
			}

			err := pq.addClientToQueue(client)
			if err != nil {
				log.Println("issue adding client to matchmaking queue:", err)
				return
			}

			log.Println("client added to queue")
		case client, ok := <- pq.removeCh:
			if !ok {
				log.Println("issue reading client from remove client channel")
				return
			}

			err := pq.removePlayerFromQueue(client)
			if err != nil {
				log.Println("issue removing client from matchmaking queue:", err)
				return
			}
		}

		currQueueLen := pq.Len()
		if currQueueLen != 0 && currQueueLen % 2 == 0 {
			players := pq.popNewMatchPlayers(2)

			// how do we link the player to the match in a more global level?
			// the maanger can maintain a map of clients -> match 
			// when the client receives an event for the match, the appropriate match can be found and the state can be updated
			// the player's manager's list of games
		} 

	}
}