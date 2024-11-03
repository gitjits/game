package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

const (
	overwhelm_threshold float64 = 5.0
)

type Unit struct {
	Name       string
	MoveRange  int
	HP         int
	StartingHP int
	Offense    int
	Defense    int
	BotUnit    bool

	// This is my replacement for an Optional<Unit> type, when going through
	// the grid we can do "if unit.Present {".
	Present bool
}

// Player units
var UNIT_PHONOMANCER Unit = Unit{Name: "Phonomancer", MoveRange: 2, HP: 4, StartingHP: 4, Offense: 6, Defense: 2, Present: true}
var UNIT_NEWTHANDS Unit = Unit{Name: "Newt-Hands", MoveRange: 3, HP: 5, StartingHP: 5, Offense: 3, Defense: 2, Present: true}
var UNIT_LBJ Unit = Unit{Name: "Lyndon B. Johnson", MoveRange: 2, HP: 8, StartingHP: 8, Offense: 3, Defense: 6, Present: true}

// Enemy units
var UNIT_WING_CENTIPEDE Unit = Unit{Name: "Tri-Winged Centipede", MoveRange: 6, HP: 3, Offense: 4, Defense: 2, Present: true, BotUnit: true}

func randomPopulate(grid *TileGrid) {
	grid.Tiles[0][1].occupant = UNIT_LBJ
	grid.Tiles[2][2].occupant = UNIT_PHONOMANCER
	grid.Tiles[1][3].occupant = UNIT_NEWTHANDS
}

func reportWinner(self *Unit, opp *Unit, g *Game) {
	var msg string
	if self.HP > 0 {
		msg = fmt.Sprintf("%s attacked %s and won!", self.Name, opp.Name)
	} else if opp.HP > 0 {
		msg = fmt.Sprintf("%s got the jump on %s and still lost!", self.Name, opp.Name)
	} else {
		msg = fmt.Sprintf("%s attacked %s, but they both lived!", self.Name, opp.Name)
	}
	g.logger.AddMessage("[!] ", msg, false)
}

func (self *Unit) attackEnemy(opp *Unit, g *Game) {
	defer reportWinner(self, opp, g)
	var imbalance float64 = float64(self.Offense - (opp.Defense))
	g.logger.AddMessage("[+] ", fmt.Sprintf("%s has a power imbalance of %f against %s", self.Name, imbalance, opp.Name), false)

	lost := math.Signbit(imbalance)

	// 20% chance of an underdog win, unless the difference is massive
	random := rand.IntN(100)
	if random > 80 && math.Abs(imbalance) < overwhelm_threshold {
		lost = !lost
	}

	if lost {
		// Attacker lost...
		self.HP = 0
	} else {
		// Attacker overcame odds!
		opp.HP = 0
	}
}
