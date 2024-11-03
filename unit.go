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
var UNIT_WING_CENTIPEDE Unit = Unit{Name: "Tri-Winged Centipede", MoveRange: 6, HP: 3, StartingHP: 3, Offense: 4, Defense: 2, Present: true, BotUnit: true}
var UNIT_WIZZY Unit = Unit{Name: "WIZZY", MoveRange: 10, HP: 1, StartingHP: 1, Offense: 8, Defense: 1, Present: true, BotUnit: true}

func randomPopulate(grid *TileGrid) {
    x := rand.IntN(9)
    y := 5 + rand.IntN(3)
    grid.Tiles[x][y].occupant = UNIT_LBJ
    x = rand.IntN(9)
    y = 5 + rand.IntN(3)
    grid.Tiles[x][y].occupant = UNIT_NEWTHANDS
    x = rand.IntN(9)
    y = 5 + rand.IntN(3)
    grid.Tiles[x][y].occupant = UNIT_PHONOMANCER

    x = rand.IntN(9)
    y = rand.IntN(3)
    grid.Tiles[x][y].occupant = UNIT_WIZZY
    x = rand.IntN(9)
    y = rand.IntN(3)
    grid.Tiles[x][y].occupant = UNIT_WIZZY
    x = rand.IntN(9)
    y = rand.IntN(3)
    grid.Tiles[x][y].occupant = UNIT_WING_CENTIPEDE
    x = rand.IntN(9)
    y = rand.IntN(3)
    grid.Tiles[x][y].occupant = UNIT_WING_CENTIPEDE
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
    if self.BotUnit == opp.BotUnit {
        g.logger.AddMessage("[!] ", fmt.Sprintf("friendly fire coming from %s!", self.Name), false)
        return
    }
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
