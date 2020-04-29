package game

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/rs/xid"
)

const CARD_COLOR_WHITE = "WHITE"
const CARD_COLOR_BLACK = "BLACK"

var whiteCards []Card
var blackCards []Card

type Player struct {
	Name   string
	Id     string
	IsCzar bool
	Cards  []Card
	Score  int
}

type Game struct {
	Id          string
	Started     bool
	Players     []Player
	BlackCard   Card
	CardsInPlay []CardInPlay
}

type Card struct {
	Id     string
	Color  string
	Text   string
	Played bool
}

type CardInPlay struct {
	Card   Card
	Player Player
}

func LoadWhiteCards() {
	file, err := os.Open("white_cards.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		if len(text) > 5 {
			whiteCards = append(whiteCards, Card{
				Color: CARD_COLOR_WHITE,
				Text:  text,
				Id:    xid.New().String(),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// ShuffleWhiteCards()
}

func ShuffleWhiteCards() {
	// Shuffle the cards
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(whiteCards), func(i, j int) { whiteCards[i], whiteCards[j] = whiteCards[j], whiteCards[i] })
}

func LoadBlackCards() {
	file, err := os.Open("black_cards.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		if len(text) > 5 {
			blackCards = append(blackCards, Card{
				Color: CARD_COLOR_BLACK,
				Text:  text,
				Id:    xid.New().String(),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	ShuffleBlackCards()
}

func ShuffleBlackCards() {
	// Shuffle the cards
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(blackCards), func(i, j int) { blackCards[i], blackCards[j] = blackCards[j], blackCards[i] })
}

func NewGame() Game {
	whiteCards = []Card{}
	blackCards = []Card{}

	LoadWhiteCards()
	LoadBlackCards()
	return Game{
		Id: xid.New().String(),
	}
}

func (g *Game) AddPlayer(p Player) Player {
	p.Id = xid.New().String()
	g.Players = append(g.Players, p)
	return p
}

func (g *Game) Start() *Game {

	ShuffleWhiteCards()

	if len(g.CardsInPlay) > 0 {
		g.Id = xid.New().String()
	}

	g.CardsInPlay = []CardInPlay{}

	rand.Seed(time.Now().UnixNano())
	czarIndex := rand.Intn(len(g.Players))

	for i, _ := range g.Players {
		g.Players[i].IsCzar = false

		if i == czarIndex {
			g.Players[i].IsCzar = true
		} else {
			newCards := DrawCards(10 - len(g.Players[i].Cards))
			for c := range newCards {
				g.Players[i].Cards = append(g.Players[i].Cards, newCards[c])
			}
		}
	}
	g.Started = true
	g.BlackCard = DrawBlackCard()
	return g
}

func (g *Game) PutCardInPlay(cardId, playerId string) {
	var cardInPlay CardInPlay
	for i := range g.Players {
		if g.Players[i].Id == playerId {
			for j := range g.Players[i].Cards {
				if g.Players[i].Cards[j].Id == cardId {
					cardInPlay.Card = g.Players[i].Cards[j]
					cardInPlay.Player = g.Players[i]

					g.CardsInPlay = append(g.CardsInPlay, cardInPlay)
					// Remove the card from player
					g.Players[i].Cards[j].Played = true
					g.Players[i].Cards = append(g.Players[i].Cards[:j], g.Players[i].Cards[j+1:]...)
					break
				}
			}
		}
	}
}

func (g *Game) AwardCardInPlay(cardId string) {
	fmt.Println(cardId)
	for i := range g.CardsInPlay {
		if g.CardsInPlay[i].Card.Id == cardId {
			playerId := g.CardsInPlay[i].Player.Id
			for j := range g.Players {
				if g.Players[j].Id == playerId {
					g.Players[j].Score = g.Players[j].Score + 1
				}
			}
		}
	}
}

func DrawCards(count int) []Card {
	var cards []Card

	for i := 1; i <= count; i++ {
		cards = append(cards, whiteCards[i])
		// Remove the card from the WhiteCards slice
		whiteCards = append(whiteCards[:i], whiteCards[i+1:]...)
	}

	fmt.Printf("Cards left: %d", len(whiteCards))
	return cards
}

func DrawBlackCard() Card {
	ShuffleBlackCards()
	return blackCards[0]
}
