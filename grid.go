package main

import (
	//"fmt"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Width    int
	Height   int
	X        int
	Y        int
	Selected bool
	Color    color.RGBA
	Radius   int
}

type TileGrid struct {
	Tiles   [][]*Tile
	X       int
	Y       int
	SizeX   int
	SizeY   int
	BoundsX int
	BoundsY int
	Color   color.RGBA
}

func drawGridTree(startTree *GridTree, screen *ebiten.Image, offsetY int) {
    tree := startTree
    offsetX := 0
    for tree != nil {
        if tree.branch != nil {
            //drawGridTree(tree, screen, offsetY + 60)
        }
	//g.grid = createGrid(0, 0, 9, 9, screenWidth/2, screenHeight/2, color.RGBA{R: 255, B: 255, G: 255, A: 1})
        tree.grid.X = offsetX
        tree.grid.Y = offsetY
        tree.grid.BoundsX = 50
        tree.grid.BoundsY = 50
        tree.grid.Update()
        drawGrid(tree.grid, screen)
        //fmt.Println("Drawing", tree.grid)
        fmt.Println(tree.grid)
        tree = tree.next
        offsetX += 60
    }
}

func drawGrid(grid TileGrid, screen *ebiten.Image) {
	for j := 0; j < len(grid.Tiles); j++ {
		for i := 0; i < len(grid.Tiles[j]); i++ {
			tile := grid.Tiles[j][i]
			if !tile.Selected {
				drawPolygon(6, tile.X+grid.X, tile.Y+grid.Y, tile.Radius, tile.Color, screen)
                fmt.Println(6, tile.X+grid.X, tile.Y+grid.Y, tile.Radius, tile.Color)
			}
		}
	}
}

func createGrid(X int, Y int, SizeX int, SizeY int, BoundsX int, BoundsY int, defaultColor color.RGBA) TileGrid {
	grid := TileGrid{
		X:       X,
		Y:       Y,
		SizeX:   SizeX,
		SizeY:   SizeY,
		BoundsX: BoundsX,
		BoundsY: BoundsY,
		Color:   defaultColor,
	}
	grid.Tiles = make([][]*Tile, grid.SizeX)
	for i := range grid.Tiles {
		grid.Tiles[i] = make([]*Tile, grid.SizeY)
	}

	for j := 0; j < grid.SizeY; j++ {
		for i := 0; i < grid.SizeX; i++ {
			grid.Tiles[j][i] = &Tile{
				Color:    defaultColor,
				Selected: false,
			}
		}
	}
	grid.Update()

	return grid
}

func (grid *TileGrid) Update() {
	for j := 0; j < grid.SizeY; j++ {
		for i := 0; i < grid.SizeX; i++ {
			r := grid.BoundsX / grid.SizeX / 2
			w, h := r*2, r*2
			ymult := 1.0
			if i%2 != 0 {
				ymult = 2
			}
			X, Y := w*i+r, h*j+int(math.Floor(float64(r)*ymult))
			grid.Tiles[j][i].Width = w
			grid.Tiles[j][i].Height = h
			grid.Tiles[j][i].X = X
			grid.Tiles[j][i].Y = Y
			grid.Tiles[j][i].Radius = r
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				mx, my := ebiten.CursorPosition()
				if mx <= X+r && mx >= X-r && my <= Y+r && my >= Y-r {
					grid.Tiles[j][i].Selected = true
				}
			}

			//fmt.Println(grid.Tiles[j][i])
		}
	}
}
