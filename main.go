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

const (
	screenWidth  = 500
	screenHeight = 500
)

type Game struct {
	selectedGrid TileGrid
	grid         TileGrid
	gridTree     *GridTree

	// Backing git repo to track grid changes
	backingFS  billy.Filesystem
	repo       *git.Repository
	cur_branch string

	// Meta stuff
	op     ebiten.DrawImageOptions
	inited bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	//g.selectedGrid = createGrid(0, 0, 9, 9, screenWidth/2, screenHeight/2, color.RGBA{R: 255, B: 255, G: 255, A: 1})
    err := gitSetup(g)
    if err != nil {
        panic(err)
    }
	g.cur_branch = "master"
	g.gridTree, err = buildCommitTree(g)
    if err != nil {
        panic(err)
    }
	//fmt.Println(g.gridTree)

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
	drawGridTree(g.gridTree, screen, 50)
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
	//ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
