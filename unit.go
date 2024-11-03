package main

import "math/rand/v2"

const (
	overwhelm_threshold int = 5
)

type Unit struct {
	Name      string
	MoveRange int
	HP        int
	Offense   int
	Defense   int

	OffenseBonus int
	DefenseBonus int

	// This is my replacement for an Optional<Unit> type, when going through
	// the grid we can do "if unit.Present {".
	Present bool
}

var UNIT_PHONOMANCER Unit = Unit{Name: "Phonomancer", MoveRange: 2, HP: 4, Offense: 6, Defense: 2, Present: true}
var UNIT_NEWTHANDS Unit = Unit{Name: "Newt-Hands", MoveRange: 3, HP: 5, Offense: 3, Defense: 2, Present: true}
var UNIT_LBJ Unit = Unit{Name: "Lyndon B. Johnson", MoveRange: 2, HP: 8, Offense: 3, Defense: 6, Present: true}

func randomPopulate(grid *TileGrid) {
	grid.Tiles[0][1].occupant = UNIT_LBJ
	grid.Tiles[2][2].occupant = UNIT_PHONOMANCER
	grid.Tiles[1][3].occupant = UNIT_NEWTHANDS
}

func (self *Unit) attackEnemy(opp *Unit) {
	var fucked_level int = self.Offense - (opp.Defense + opp.DefenseBonus)
	if fucked_level > overwhelm_threshold {
		// Overwhelming power difference, attacker instantly wins
		opp.HP = 0
		return
	} else if fucked_level < -overwhelm_threshold {
		// Overwhelming power difference in favor of defender, attacker loses
		self.HP = 0
		return
	}

	// The difference in power isn't massive, it could go either way
	var winner = rand.IntN(fucked_level) - (fucked_level / 2)
	if winner > 0 {
		opp.HP = 0
	} else {
		self.HP = 0
	}
}
