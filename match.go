package main

// match info and match state info
type Match struct {
	players []*Client
}

func CreateMatch(players []*Client) (Match, error){
	//create match object
	// adds players to the match
	// notifies users that they found a match
	//notifiying users that a match was found would involve sending an 
	//html tempalte	
	// should this launch a goroutine?
}