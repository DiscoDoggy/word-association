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

	playerExitCh chan *Client

	usedWords	map[string]bool

	endGameCh chan bool
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

		wordSubCh: make(chan string),
		wordEditCh: make(chan string),

		endGameCh: make(chan bool),

		playerExitCh: make(chan *Client),
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
	// how to handle exit game and game end clean up
		// just make an event that gets sent on button exit??
		//page redirections are really just pushes of html so in this case we would push the home page to the user again after cleaning up
	for {
		select {
		case word, ok := <- m.wordSubCh:
			if !ok {
				log.Println("issue reading word from word submission channel")
			}
			//TODO: Preprocess word. eg Mcdonald's === McDonalds for say the fast food topic
			log.Println("entered m.WordSUbmission channel event")
			log.Println("entered word:", word)
			_, doesExist := topic[word]
			_, isAlreadyUsedWord := m.usedWords[word]
			if !doesExist || isAlreadyUsedWord {
				//send lose game and win game to respective clients
				err := m.SendEndGameHtmlToPlayers(currPlayerTurn)
				if err != nil {
					log.Println("error sending end game state to players")
					return
				}
				log.Println("enter already used word || word does not exist in topic")
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
		case exitingPlayer, ok := <- m.playerExitCh:
			if !ok {
				log.Println("error reading player leave channel")
			} else {
				for i, player := range m.players {
					if exitingPlayer == player {
						m.players = append(m.players[:i], m.players[i + 1:]...) 
					}
				}	

				html, err := templates.ConvertComponentToHtml(templates.WrappedIndex(templates.HomeContent()))
				if err != nil {
					log.Println("error converting home page template to html")
					return
				}
				//we want to wait until all players have formally exited before we destory the match
				exitingPlayer.egress <- html
				if len(m.players) <= 0 {
					m.matchManager.removeMatchCh <- m	
					return
				}
			}

		}
	}
}

func (m *Match) SendEndGameHtmlToPlayers(losingPlayer *Client) error {
	html, err := templates.ConvertComponentToHtml(templates.EndgameCard(false, ""))
	if err != nil {
		//TODO: Create error banner
		log.Println("error rendering end game cards:", err)
		return err
	}

	losingPlayer.egress <- html

	//TODO later determine how to determine singular winner if there were more than 2 players
	html, err = templates.ConvertComponentToHtml(templates.EndgameCard(true, ""))
	if err != nil {
		log.Println("error rendering end game cards:", err)
	}

	for _, player := range m.players {
		if player != losingPlayer {
			player.egress <- html
		}
	}

	return nil
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

