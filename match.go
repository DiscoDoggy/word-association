package main

import (
	"fmt"
	"log"
	"math/rand"
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

	wordSubCh chan string
	wordEditCh chan string

	usedWords	map[string]bool
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
		addMatchCh:    make(chan *Match),
		removeMatchCh: make(chan *Match),
	}
}

func CreateMatch(players []*Client) {

	fmt.Println("Creating match")
	match := Match{
		id:           uuid.New(),
		players:      players,
		matchManager: players[0].manager.matchManager,
		isActive:     false,
		usedWords: make(map[string]bool),
	}

	log.Println("entering match manager adding")

	players[0].manager.matchManager.addMatchCh <- &match

	log.Println("completed adding match to manager")

	//send both clients the htmx then play game goroutine??

	// for getting and checking words, we can preload into memory or inside of the match object the list of valid words from the database
	// this would avoid making a database call every time we need to check a word and would significantly increase speed of word checks

	go match.PlayGame()
}

func (m *Match) pickRandomTopic() map[string]bool {
	randIndex := rand.Intn(len(topics))

	var topicName string
	counter := 0

	for key, _ := range topics {
		if counter == randIndex {
			topicName = key
			break	
		}

		counter++
	}

	topic := topics[topicName]	

	log.Println("topic name:", topicName)
	return topic
}

func (m *Match) PlayGame() {
	//need to utilize some timers here
	// introduce the topic on screen somewhere for 2 seconds
	// count down the game with big numbers on screen 
	// then play game loop

	//choose topic
	topic := m.pickRandomTopic() // this should proabbly go inside of the CreateMatch

	//pick which player goes first
	currPlayerIdx := rand.Intn(len(m.players))
	currPlayerTurn := m.players[currPlayerIdx]

	err := m.SendGameHtmlToCorrectPlayer(currPlayerTurn)
	if err != nil {
		log.Println("error sending game html to appropriate players in match initialization:", err)
		return
	}

	log.Println("html for game page sent")

	//game loop
	for {
		select {
		case word, ok := <- m.wordSubCh:
			if !ok {
				log.Println("issue reading word from word submission channel")
			}
			//TODO: Preprocess word. eg Mcdonald's === McDonalds for say the fast food topic
			_, doesExist := topic[word]
			_, isAlreadyUsedWord := m.usedWords[word]
			if !doesExist || isAlreadyUsedWord {
				//send lose game and win game to respective clients
			} else {
				m.usedWords[word] = true

				if currPlayerIdx + 1 >= len(m.players) {
					currPlayerIdx = 0
				} else {
					currPlayerIdx += 1
				}

				currPlayerTurn = m.players[currPlayerIdx]
				// send the appropriate htmx to clients	
				err = m.SendGameHtmlToCorrectPlayer(currPlayerTurn)
				if err != nil {
					log.Println("error sending game html to appropriate players in match initialization:", err)
					return
				}
			}
		}
	}
}

func (m *Match) SendGameHtmlToCorrectPlayer (currPlayerTurn *Client) error {
	gameContentHtml, err := templates.ConvertComponentToHtml(templates.NewGameContent(true))
	if err != nil {
		//TODO: Create error banner htmx and send to clients
		log.Println("error rendering game page:", err)
		return err
	}
	
	currPlayerTurn.egress <- gameContentHtml

	for _, player := range m.players {
		if player != currPlayerTurn {
			gameContentHtml, err = templates.ConvertComponentToHtml(templates.NewGameContent(false))
			if err != nil {
				//TODO: Create error banner htmx and send to clients
				log.Println("error rendering game page:", err)
				return err
			}

			player.egress <- gameContentHtml
		}
	}

	return nil
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

			log.Println("doing match manager add match")

			matchManager.matches = append(matchManager.matches, match)
			for _, client := range match.players {
				matchManager.clientToMatch[client] = match
			}

			log.Println("added match to manager...returning")

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
					break
				}
			}
		}
	}
}

// who should launch the game and what should launch the game what event
