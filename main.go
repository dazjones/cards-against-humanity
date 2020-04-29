package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dazjones/cards-against-humanity/game"
)

var g game.Game

func renderJSON(w http.ResponseWriter, g interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	g = game.NewGame()
	renderJSON(w, g)
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	g.Start()
	renderJSON(w, g)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, g)
}

func cardPlayHandler(w http.ResponseWriter, r *http.Request) {
	cardIdRaw := fmt.Sprintf("%s", r.URL.Query()["card"])
	str1 := strings.Replace(cardIdRaw, "[", "", -1)
	cardId := strings.Replace(str1, "]", "", -1)

	playerIdRaw := fmt.Sprintf("%s", r.URL.Query()["player"])
	str2 := strings.Replace(playerIdRaw, "[", "", -1)
	playerId := strings.Replace(str2, "]", "", -1)

	g.PutCardInPlay(cardId, playerId)
	renderJSON(w, g)
}

func cardAwardHandler(w http.ResponseWriter, r *http.Request) {
	cardIdRaw := fmt.Sprintf("%s", r.URL.Query()["card"])
	str1 := strings.Replace(cardIdRaw, "[", "", -1)
	cardId := strings.Replace(str1, "]", "", -1)

	g.AwardCardInPlay(cardId)
	renderJSON(w, g)
}

func newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	name := fmt.Sprintf("%s", r.URL.Query()["name"])
	var p game.Player

	if name != "[]" {
		p = g.AddPlayer(game.Player{
			Name: name,
		})
	}

	renderJSON(w, p)
}

func main() {
	http.HandleFunc("/api/game/cards/award", cardAwardHandler)
	http.HandleFunc("/api/game/cards/play", cardPlayHandler)
	http.HandleFunc("/api/game/new", newGameHandler)
	http.HandleFunc("/api/game/players/new", newPlayerHandler)
	http.HandleFunc("/api/game/start", startGameHandler)
	http.HandleFunc("/api/game", gameHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
