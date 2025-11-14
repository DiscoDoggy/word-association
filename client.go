package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id			uuid.UUID	
	connection *websocket.Conn
	manager    *Manager

	//to avoid concurrent writes to the websocket
	egress chan []byte
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		id: uuid.New(),
		connection: conn,
		manager: manager,
		egress: make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func () {
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
			if eventType != "" {
				// raw json strings come with their leading and trailing quotation marks still this parses this out
				eventType = eventType[1:len(eventType) - 1]
			}

			clientEventType, err := GetClientEventFromStr(eventType)
			if err != nil {
				log.Fatal(err)
			}
			if clientEventType == mmQueue {
				err = c.manager.playerQueue.addClientToQueue(c)
				if err != nil {
					log.Fatal(err)
				}
			} else if clientEventType == mmUnqueue {
				err = c.manager.playerQueue.removePlayerFromQueue(c)
				if err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println("Client event:", string(event["event"]))
			fmt.Println("Client event message:", string(event["username-input"]))
			
		}

		//for testing
		for wsclient := range c.manager.clients {
			wsclient.egress <- payload
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
			}
			log.Println("sent message")
		}
	}
}

//so one of the difficulties here is that t

//write messages sends  messages tot he clients
// read messages is messages received by the client object on the server and can be readtjerltjwerkltwtwtwetwertw