package main

import (
	"bytes"
	"fmt"
	"math"

	//"math"

	//"fmt"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"

	//"math"
	//"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 500
	screenHeight = 500
	maxAngle     = 256
    gridRadius = 25
)

var (
	ebitenImage *ebiten.Image
)

//go:embed hex.png
var TileIMG []byte

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(TileIMG))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)
	ebitenImage = ebiten.NewImage(64, 64)

	op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(0.25, 0.25)
	op.ColorScale.ScaleAlpha(0.5)
	ebitenImage.DrawImage(origEbitenImage, op)
}

type Tile struct {
	width int
	height int
	x int
	y int
    lx int
    ly int
}

type TileGrid struct {
	tiles [][]*Tile
    x int
    y int
}

func (s *TileGrid) Update() {
}

type Game struct {
	touchIDs []ebiten.TouchID
	grid     TileGrid
	op       ebiten.DrawImageOptions
	inited   bool
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

    g.grid.x = 10
    g.grid.y = 10
	g.grid.tiles = make([][]*Tile, g.grid.x)
    for i := range g.grid.tiles {
        g.grid.tiles[i] = make([]*Tile, g.grid.y)
    }
    for j := 0; j < g.grid.y; j++ {
        for i := 0; i < g.grid.x; i++ {
            w, h := gridRadius * 2 - 10, gridRadius * 2 - 3
            ymult := 1.0
            if i % 2 != 0 {
                ymult = 2
            }
            x, y := w * i + gridRadius, h * j + int(math.Floor(gridRadius*ymult))
            g.grid.tiles[j][i] = &Tile{
                width:  w,
                height: h,
                x: x,
                y: y,
                lx: i,
                ly: j,
            }
            fmt.Println(g.grid.tiles[j][i])
        }
    }
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	g.grid.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x4C, 0x4C, 0xFF});
	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	//w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	for j := 0; j < len(g.grid.tiles); j++ {
        for i := 0; i < len(g.grid.tiles[j]); i++ {
            s := g.grid.tiles[j][i]
            //g.op.GeoM.Reset()
            //g.op.GeoM.Rotate(math.Pi/4)
            //g.op.GeoM.Translate(float64(s.x), float64(s.y))
            //screen.DrawImage(ebitenImage, &g.op)
            drawPolygon(6, s.x, s.y, gridRadius, color.RGBA{R: 255, B: 255, G: 255}, screen)
        }
	}
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
