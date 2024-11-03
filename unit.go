package main

import "math/rand/v2"

const (
	overwhelm_threshold int = 5
)

type Unit struct {
	move_range int
	hp         int
	offense    int
	defense    int

	offense_bonus int
	defense_bonus int
}

func (self *Unit) attackEnemy(opp *Unit) {
	var fucked_level int = self.offense - (opp.defense + opp.defense_bonus)
	if fucked_level > overwhelm_threshold {
		// Overwhelming power difference, attacker instantly wins
		opp.hp = 0
		return
	} else if fucked_level < -overwhelm_threshold {
		// Overwhelming power difference in favor of defender, attacker loses
		self.hp = 0
		return
	}

	// The difference in power isn't massive, it could go either way
	var winner = rand.IntN(fucked_level) - (fucked_level / 2)
	if winner > 0 {
		opp.hp = 0
	} else {
		self.hp = 0
	}
}
