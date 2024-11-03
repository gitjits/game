package main

import (
	//"fmt"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Selected bool
	Color    color.RGBA
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
    IsSelectedGrid bool
}

func drawGridTree(startTree *GridTree, screen *ebiten.Image, offsetY int) {
	tree := startTree
	offsetX := 0
	for tree != nil {
		if tree.branch != nil {
			//drawGridTree(tree, screen, offsetY + 60)
		}
        if !tree.grid.IsSelectedGrid {
            tree.grid.X = offsetX
            tree.grid.Y = offsetY
            tree.grid.BoundsX = 50
            tree.grid.BoundsY = 50
        } else {
            tree.grid.X = screenWidth/4
            tree.grid.Y = screenHeight/4
            tree.grid.BoundsX = screenWidth/2
            tree.grid.BoundsY = screenHeight/2
            faux := createGrid(offsetX, offsetY, tree.grid.SizeX, tree.grid.SizeY, 50, 50, tree.grid.Color)
            faux.Tiles = tree.grid.Tiles
            drawGrid(faux, screen)
        }
		tree.grid.Update()
		drawGrid(tree.grid, screen)
		//fmt.Println(tree.grid)
		tree = tree.next
		offsetX += 60
	}
}

func drawGrid(grid TileGrid, screen *ebiten.Image) {
	for j := 0; j < len(grid.Tiles); j++ {
		for i := 0; i < len(grid.Tiles[j]); i++ {
			tile := grid.Tiles[j][i]
			if !tile.Selected {
				Xpos, Ypos := tileScreenPos(&grid, i, j)
				r := tileRadius(&grid)
				drawPolygon(6, Xpos, Ypos, r, tile.Color, screen)
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

func tileRadius(grid *TileGrid) int {
	return grid.BoundsX / grid.SizeX / 2
}

func tileScreenPos(grid *TileGrid, row int, col int) (int, int) {
	r := tileRadius(grid)
	w, h := r*2, r*2
	ymult := 1.0
	if (row % 2) != 0 {
		ymult = 2
	}

	var X = grid.X + w*row + r
	var Y = grid.Y + h*col + int(math.Floor(float64(r)*ymult))

	return X, Y
}

func (grid *TileGrid) Update() {
    if grid.IsSelectedGrid {
        if ebiten.IsKeyPressed(ebiten.KeyEscape) {
            grid.IsSelectedGrid = false
            fmt.Println("KILLED")
        }
        for j := 0; j < grid.SizeY; j++ {
            for i := 0; i < grid.SizeX; i++ {
                X, Y := tileScreenPos(grid, i, j)
                r := tileRadius(grid)
                if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
                    mx, my := ebiten.CursorPosition()
                    if mx <= X+r && mx >= X-r && my <= Y+r && my >= Y-r {
                        grid.Tiles[j][i].Selected = true
                    }
                }

                //fmt.Println(grid.Tiles[j][i])
            }
        }
    } else {
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            mx, my := ebiten.CursorPosition()
            if mx >= grid.X && mx <= grid.X+grid.BoundsX && my <= grid.Y+grid.BoundsY && my >= grid.Y {
                grid.IsSelectedGrid = true
            }
        }
    }
}
