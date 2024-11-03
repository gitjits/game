package main

import (
	//"fmt"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func drawGridTree(tree *GridTree, screen *ebiten.Image) {
	// Draw current node
	if tree.grid.SizeX == 0 {
		fmt.Println("GRID IS NIL", tree.grid)
		return
	}
	var curX int = screenWidth
	cur_tree := *tree
	for cur_tree.prev != nil {
		curX -= 60
		// Handle selected grid
		if cur_tree.grid.IsSelectedGrid {
			// Draw main selected grid in center
			cur_tree.grid.X = screenWidth / 4
			cur_tree.grid.Y = screenHeight / 4
			cur_tree.grid.BoundsX = screenWidth / 2
			cur_tree.grid.BoundsY = screenHeight / 2
			ebitenutil.DebugPrintAt(screen, cur_tree.commitHash, screenWidth/2, screenHeight*(3/4))

			// Draw small version in tree
			faux := createGrid(curX, screenHeight/2, cur_tree.grid.SizeX, cur_tree.grid.SizeY, 50, 50, cur_tree.grid.Color)
			faux.Tiles = cur_tree.grid.Tiles
			drawGrid(faux, screen)
		} else {
			cur_tree.grid.X = curX
			cur_tree.grid.Y = (cur_tree.generation * 60) + (screenHeight / 2)
			cur_tree.grid.BoundsX = 50
			cur_tree.grid.BoundsY = 50
		}

		cur_tree.grid.Update()
		drawGrid(cur_tree.grid, screen)

		fmt.Printf("prev node %p, next %p, gen %d, %dx%d\n", cur_tree.prev, cur_tree.next, cur_tree.generation, cur_tree.grid.SizeX, cur_tree.grid.SizeY)
		cur_tree = *cur_tree.prev
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
