package game_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dazjones/cards-against-humanity/game"
)

var _ = Describe("Game", func() {
	Describe("New", func() {
		Context("New Game", func() {
			g := NewGame()

			It("Should have an ID", func() {
				Expect(g.Id).NotTo(Equal(""))
			})

			It("Should have a unique ID", func() {
				Expect(NewGame().Id).NotTo(Equal(NewGame().Id))
			})
			It("Should have have no players", func() {
				Expect(len(g.Players)).To(Equal(0))
			})
		})

		Context("Add Players", func() {
			g := NewGame()

			It("Should have one player", func() {
				g.AddPlayer(Player{Name: "Darren"})
				Expect(len(g.Players)).To(Equal(1))
			})

			It("Should have zero cards for each player", func() {
				Expect(len(g.Players[0].Cards)).To(Equal(0))
			})

			It("Should have two player", func() {
				g.AddPlayer(Player{Name: "Robert"})
				Expect(len(g.Players)).To(Equal(2))
			})
		})
	})
	Describe("Start", func() {
		g := NewGame()

		g.AddPlayer(Player{Name: "Darren"})
		g.AddPlayer(Player{Name: "Robert"})

		g.Start()

		It("each player draws ten white cards", func() {
			for _, player := range g.Players {
				Expect(len(player.Cards)).To(Equal(10))
			}

			for _, player := range g.Players {
				for _, card := range player.Cards {
					Expect(card.Color).To(Equal(CARD_COLOR_WHITE))
				}
			}
		})

		It("each player should have unique cards", func() {
			Expect(g.Players[0].Cards[0]).NotTo(Equal(g.Players[1].Cards[0]))
			Expect(g.Players[0].Cards[2]).NotTo(Equal(g.Players[1].Cards[2]))
			Expect(g.Players[0].Cards[4]).NotTo(Equal(g.Players[1].Cards[4]))
			Expect(g.Players[0].Cards[6]).NotTo(Equal(g.Players[1].Cards[6]))
			Expect(g.Players[0].Cards[8]).NotTo(Equal(g.Players[1].Cards[8]))
		})
	})
})
