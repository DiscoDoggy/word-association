package main

import (
	"context"
	"log"
	"net/http"

	"github.com/DiscoDoggy/word-association/templates"
	// "github.com/a-h/templ"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "method not found", http.StatusNotFound)
		return
	}

	homePage := templates.HomePage()
	homePage.Render(context.Background(), w)
}

// func servePartyRoom(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/party-room" {
// 		http.Error(w, "not found", http.StatusNotFound)
// 		return
// 	}

// 	if r.Method != "GET" {
// 		http.Error(w, "method not found", http.StatusNotFound)
// 		return
// 	}

// 	partyPage := templates.PartyPage()
// 	partyPage.Render(context.Background(), w)
// }

// func serveMatchRoom(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/match" {
// 		http.Error(w, "not found", http.StatusNotFound)
// 		return
// 	}

// 	matchPage := templates.GamePage() 
// 	matchPage.Render(context.Background(), w)
// }

func main() {
	wsManager := NewManager()


	http.HandleFunc("/", serveIndex)
	// http.HandleFunc("/party-room", servePartyRoom)
	// http.HandleFunc("/match", serveMatchRoom)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/ws", wsManager.serveWS)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
