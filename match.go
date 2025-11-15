package main

import "fmt"

// match info and match state info
type Match struct {
	players []*Client
}

func CreateMatch(players []*Client) error {

	fmt.Println("Creating match")
	//create match object
	// adds players to the match
	// notifies users that they found a match
	//notifiying users that a match was found would involve sending an 
	//html tempalte	
	// should this launch a goroutine?
	// it probably should so it can manage itself and have its own game loop
	return nil
}