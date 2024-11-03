package main

import (
	//"fmt"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tile struct {
	Selected bool
	Color    color.RGBA
}

type TileGrid struct {
	Tiles          [][]*Tile
	X              int
	Y              int
	SizeX          int
	SizeY          int
	BoundsX        int
	BoundsY        int
	Color          color.RGBA
	IsSelectedGrid bool
}

func drawGridTree(g *Game, tree *GridTree, screen *ebiten.Image, offsetY, offsetX int) {
	// Draw current node
	if tree.grid.SizeX == 0 {
		fmt.Println("GRID IS NIL", tree.grid)
		return
	}

	// Handle selected grid
    if g.selected != nil {
        fmt.Println(g.selected.Color)
		g.selected.X = screenWidth / 2 - 150
		g.selected.Y = screenHeight / 2 - 150
		g.selected.BoundsX = 310
		g.selected.BoundsY = 310
        g.selected.Update(g)
        if g.selected != nil {
            drawGrid(*g.selected, screen)
        }
    }
	if tree.grid.IsSelectedGrid {
		// Draw main selected grid in center

		// Draw small version in tree
		faux := createGrid(offsetX, offsetY + tree.generation*120, tree.grid.SizeX, tree.grid.SizeY, 110, 110, tree.grid.Color)
		faux.Tiles = tree.grid.Tiles
		drawGrid(faux, screen)
	} else {
		tree.grid.X = offsetX
		tree.grid.Y = offsetY + tree.generation*120
		tree.grid.BoundsX = 110
		tree.grid.BoundsY = 110
	}

	tree.grid.Update(g)
	drawGrid(tree.grid, screen)

	// Continue main branch
	if tree.next != nil && tree.next.grid.SizeX != 0 {
		drawGridTree(g, tree.next, screen, offsetY, offsetX+120)
	}
}

func drawGrid(grid TileGrid, screen *ebiten.Image) {
	r := tileRadius(&grid)
	for j := 0; j < len(grid.Tiles); j++ {
		for i := 0; i < len(grid.Tiles[j]); i++ {
			tile := grid.Tiles[j][i]
			if !tile.Selected {
				Xpos, Ypos := tileScreenPos(&grid, i, j)
				drawPolygon(6, Xpos, Ypos, r, tile.Color, screen)
			}
		}
	}
	vector.StrokeRect(screen, float32(grid.X-r/2), float32(grid.Y), float32(grid.BoundsX+r), float32(grid.BoundsY+r), 1, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)
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
            if i % 2 == 0 {
                grid.Tiles[j][i] = &Tile{
                    Color:    defaultColor,
                    Selected: false,
                }
            } else {
                grid.Tiles[j][i] = &Tile{
                    Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
                    Selected: false,
                }
            }
		}
	}

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

func (grid *TileGrid) Update(g *Game) {
	if grid.IsSelectedGrid {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			grid.IsSelectedGrid = false
            g.selected = nil
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
			}
		}
	} else {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if mx >= grid.X && mx <= grid.X+grid.BoundsX && my <= grid.Y+grid.BoundsY && my >= grid.Y {
				grid.IsSelectedGrid = true
                if g.selected != nil {
                    g.selected.IsSelectedGrid = false
                }
                g.selected = grid
			}
		}
	}
}
