package main

import (
	"fmt"
	"log"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	screenWidth  = 1000
	screenHeight = 700
)

type Game struct {
	grid     TileGrid
	gridTree GridTree
	selected *TileGrid

	// Meta stuff
	op     ebiten.DrawImageOptions
	inited bool

	logger              *LogWindow
	scrollX             int
	autoScroll          bool
	CPressedLastFrame   bool
	BPressedLastFrame   bool
	MPressedLastFrame   bool
	RPressedLastFrame   bool
	ESCPressedLastFrame bool

	infoSprite    Unit
	hidden        bool
	stop          bool
	botWaitPeriod int
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()
	g.scrollX = 50
	g.botWaitPeriod = -1

	err := loadEmbeddedImage()
	if err != nil {
		panic(err)
	}

	g.gridTree = GridTree{}

	g.logger = NewLogWindow()

	// Create basic test data in the repo
	g.grid = createGrid(0, 0, 9, 9, screenWidth/2, screenHeight/2, color.RGBA{R: 255, B: 255, G: 255, A: 200})
	g.grid.Update(g)
	gitCommitGrid(g, g.grid, false, true)
	g.selected = &g.grid
	g.selected.IsSelectedGrid = true

	g.grid = createGrid(0, 0, 9, 9, screenWidth/2, screenHeight/2, color.RGBA{R: 255, B: 255, G: 255, A: 200})
	g.grid.Update(g)
	gitCommitGrid(g, g.grid, false, true)
	g.selected = &g.grid
	g.selected.IsSelectedGrid = true
	randomPopulate(g.selected)
	//commitTestData(g)

	// We need a simplified commit tree to efficiently render it
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("", "", false)
	g.logger.AddMessage("ur_enemy$ ", "git init", false)
	g.logger.AddMessage("ur_enemy$ ", "git commit -m 'welcome to the game'", false)
	g.logger.AddMessage("", fmt.Sprintf("[main %s] welcome to the game", g.gridTree.commitHash[0:8]), false)
	g.logger.AddMessage("", "1 files changed, 1 insertions(+), 0 deletions(-)", false)
	g.logger.AddMessage("", "create mode 100644 board.bson", false)
	g.logger.AddMessage("you$ ", "git --help", false)
	g.logger.AddMessage("[!] ", "Controls", false)
	g.logger.AddMessage("[!] ", "esc: hide board", false)
	g.logger.AddMessage("[!] ", "c: commit", false)
	g.logger.AddMessage("[!] ", "b: new branch", false)
	g.logger.AddMessage("[!] ", "m: merge", false)
	g.logger.AddMessage("[!] ", "r: revert", false)
	fmt.Print("Setup Git repo!\n")
}

func (g *Game) Update() error {
	if g.stop {
		return nil
	}
	if !g.inited {
		g.init()
	}
	if g.botWaitPeriod == 0 {
		// Make bot move
		g.selected.makeBotMove(g)
		g.botWaitPeriod = -1
	} else {
		g.botWaitPeriod--
	}

	CPressedNow := ebiten.IsKeyPressed(ebiten.KeyC)
	if CPressedNow && !g.CPressedLastFrame {
		newGrid := g.gridTree.grid.Clone()
		gitCommitGrid(g, newGrid, false, false)
	}
	g.CPressedLastFrame = CPressedNow
	BPressedNow := ebiten.IsKeyPressed(ebiten.KeyB)
	if BPressedNow && !g.BPressedLastFrame {
		newGrid := g.gridTree.grid.Clone()
		gitCommitGrid(g, newGrid, true, false)
	}
	g.BPressedLastFrame = BPressedNow
	MPressedNow := ebiten.IsKeyPressed(ebiten.KeyM)
	if MPressedNow && !g.MPressedLastFrame {
		mergeCurrentBranch(g)
	}
	g.MPressedLastFrame = MPressedNow
	RPressedNow := ebiten.IsKeyPressed(ebiten.KeyR)
	if RPressedNow && !g.RPressedLastFrame {
		nukeCurrentBranch(g)
	}
	g.RPressedLastFrame = RPressedNow
	ESCPressedNow := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if ESCPressedNow && !g.ESCPressedLastFrame {
		g.hidden = !g.hidden
	}
	g.ESCPressedLastFrame = ESCPressedNow

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x4C, 0x4C, 0xFF})
	t2 := g.gridTree
	for t2.prev != nil && t2.prev.grid.SizeX != 0 {
		t2 = *t2.prev
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if t2.next.grid.X < 185 {
			g.scrollX += 1
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || g.autoScroll {
		if g.gridTree.grid.X > 840 {
			g.scrollX -= 1
		} else {
			g.autoScroll = false
		}
	}
	if g.infoSprite.Present {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s:\n\tHP: %d/%d\n\tDefense: %d\n\tOffense: %d", g.infoSprite.Name, g.infoSprite.HP, g.infoSprite.StartingHP, g.infoSprite.Defense, g.infoSprite.Offense), screenWidth-110, 0)
	}
	drawGridTree(g, &t2, screen, 50, g.scrollX)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
	g.logger.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Sprites (Ebitengine Demo)")
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
