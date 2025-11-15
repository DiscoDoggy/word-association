package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/DiscoDoggy/word-association/templates"
	"github.com/google/uuid"
)

// match info and match state info
type Match struct {
	id           uuid.UUID
	players      []*Client
	matchManager *MatchManager
	isActive     bool
}

type MatchManager struct {
	matches       []*Match
	clientToMatch map[*Client]*Match

	addMatchCh    chan *Match
	removeMatchCh chan *Match

	sync.RWMutex
}

func CreateMatchManager() *MatchManager {
	return &MatchManager{
		matches:       make([]*Match, 0),
		clientToMatch: make(map[*Client]*Match),
	}
}

func CreateMatch(players []*Client) {

	fmt.Println("Creating match")
	// match := Match{
	// 	id:           uuid.New(),
	// 	matchManager: players[0].manager.matchManager,
	// 	isActive:     false,
	// }

	// players[0].manager.matchManager.addMatchCh <- &match

	//send both clients the htmx then play game goroutine??
	var b bytes.Buffer
	err := templates.GameContent().Render(context.Background(), &b)
	if err != nil {
		//TODO: Create error banner htmx and send to clients
		log.Println("error rendering game page:", err)

		return
	}

	gameContentHtml := b.Bytes()

	for _, player := range players {
		player.egress <- gameContentHtml
	}

	log.Println("html for game page sent")

}

func (m *Match) PlayGame() {
	return
}

func (matchManager *MatchManager) LaunchMatchManager() {
	defer func() {
		matchManager.matches = nil
		matchManager.clientToMatch = make(map[*Client]*Match)
	}()
	for {
		select {
		case match, ok := <-matchManager.addMatchCh:
			if !ok {
				log.Println("error getting match from add channel")
				return
			}

			matchManager.matches = append(matchManager.matches, match)
			for _, client := range match.players {
				matchManager.clientToMatch[client] = match
			}

			return
		case match, ok := <-matchManager.removeMatchCh:
			if !ok {
				log.Println("error getting match from remove channel")
				return
			}

			for _, player := range match.players {
				delete(matchManager.clientToMatch, player)
			}

			for i, managedMatch := range matchManager.matches {
				if managedMatch == match {
					matchManager.matches = append(matchManager.matches[:i], matchManager.matches[i+1:]...)
				}
			}
		}
	}
}

// who should launch the game and what should launch the game what event
