package main

import (
	"log"
	"net/http"
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

	http.ServeFile(w, r, "templates/pages/home-page/index.html")
}

func servePartyRoom(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/party-room" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "method not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "templates/pages/party-room/party.html")
}

func serveMatchRoom(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/match" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "templates/pages/game-page/game.html")
}

func main() {
	wsManager := NewManager()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/party-room", servePartyRoom)
	http.HandleFunc("/match", serveMatchRoom)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	http.HandleFunc("/ws", wsManager.serveWS)

	log.Fatal(http.ListenAndServe(":3000", nil))
} 