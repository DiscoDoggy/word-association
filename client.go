package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id         uuid.UUID
	connection *websocket.Conn
	manager    *Manager

	//to avoid concurrent writes to the websocket
	egress chan []byte
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		id:         uuid.New(),
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		log.Println("Message type: ", messageType)
		log.Println("Payload: ", string(payload))

		// parse payload
		if messageType == 1 {
			var event map[string]json.RawMessage
			err := json.Unmarshal(payload, &event)
			if err != nil {
				log.Fatal(err)
			}
			// TODO : Create more robust handling for events that may not exist or may not match structure
			// fatalling inside of a goroutine kills entire program and probably should not be used
			eventType := string(event["event"])
			log.Println("event type:", eventType)
			if eventType != "" {
				// raw json strings come with their leading and trailing quotation marks still this parses this out
				eventType = TrimFirstLastChar(eventType) 
			}

			clientEventType, err := GetClientEventFromStr(eventType)
			if err != nil {
				log.Fatal(err)
			}
			if clientEventType == mmQueue {
				c.manager.playerQueue.addCh <- c
			} else if clientEventType == mmUnqueue {
				c.manager.playerQueue.removeCh <- c
			} else if clientEventType == mWordSubmission {
				log.Println("attempting to write to match word submission channel")
				match, ok := c.manager.matchManager.clientToMatch[c]
				if !ok {
					log.Println("Error: Client does not exist in client->match table")
				} else {
					submittedWord := string(event["word"])
					submittedWord = TrimFirstLastChar(submittedWord)
					match.wordSubCh <- submittedWord 
				}
			} else if clientEventType == mPlayerExit {
				match, ok := c.manager.matchManager.clientToMatch[c]
				if !ok {
					log.Println("client is not part of any match")
					//TODO: Some sort of error handling
				} else {
					match.playerExitCh <- c
				}
			}

			fmt.Println("Client event:", string(event["event"]))
			fmt.Println("Client event message:", string(event["username-input"]))
			fmt.Println("Player queue length:", len(c.manager.playerQueue.queue))
		}

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}

				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
				return
			}
			log.Println("sent message")
		}
	}
}

func TrimFirstLastChar(s string) string {
	return s[1:len(s) - 1]
}

//so one of the difficulties here is that t

//write messages sends  messages tot he clients
// read messages is messages received by the client object on the server and can be readtjerltjwerkltwtwtwetwertw
