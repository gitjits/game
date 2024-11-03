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

var UNIT_PHONOMANCER Unit = Unit{Name: "Phonomancer", MoveRange: 3, HP: 5, Offense: 6, Defense: 2, Present: true}

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
