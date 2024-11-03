package main

import (
	"fmt"

	_ "embed"
	"image/color"
	_ "image/png"
	"log"

	billy "github.com/go-git/go-billy/v5"
	git "github.com/go-git/go-git/v5"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screenWidth  = 1000
	screenHeight = 700
)

type Game struct {
	grid     TileGrid
	gridTree GridTree
    selected *TileGrid

	// Backing git repo to track grid changes
	backingFS  billy.Filesystem
	repo       *git.Repository
	cur_branch string

	// Meta stuff
	op     ebiten.DrawImageOptions
	inited bool

    logger *LogWindow
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	g.gridTree = GridTree{}

    g.logger = NewLogWindow()

	// Create basic test data in the repo
	g.grid = createGrid(0, 0, 5, 5, screenWidth/2, screenHeight/2, color.RGBA{R: 255, B: 255, G: 255, A: 1})
	gitCommitGrid(g, g.grid, false)
	commitTestData(g)

	// We need a simplified commit tree to efficiently render it
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")
    g.logger.AddMessage("ena")

	fmt.Print("Setup Git repo!\n")
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x4C, 0x4C, 0xFF})
    t2 := g.gridTree
    for t2.prev != nil && t2.prev.grid.SizeX != 0 {
        t2 = *t2.prev
    }
	drawGridTree(g, &t2, screen, 50, 50)
    g.logger.Draw(screen)
	//drawGrid(g.grid, screen)
	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	//w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
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
