package main

import (
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

func (self *Unit) attackEnemy(opp *Unit) {
	var imbalance float64 = float64(self.Offense - (opp.Defense))
	if imbalance > overwhelm_threshold {
		// Overwhelming power difference, attacker instantly wins
		opp.HP = 0
		return
	} else if imbalance < -overwhelm_threshold {
		// Overwhelming power difference in favor of defender, attacker loses
		self.HP = 0
		return
	}

	// The difference in power isn't massive, it could go either way.
	// Any random result > 1x will result in a win, which is possible but
	// less likely.
	imbalanceMag := math.Abs(imbalance)
	random := rand.IntN(int(imbalanceMag * 1.3))
	if imbalanceMag-float64(random) < 0 {
		// Attacker overcame odds!
		opp.HP = 0
	} else {
		// Attacker lost...
		self.HP = 0
	}
}
