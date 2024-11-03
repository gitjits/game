package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tile struct {
	Selected bool
	Color    color.RGBA
	occupant Unit
}

type vec2i struct {
	x     int
	y     int
	valid bool // This is only set on once the position is initialized
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
    ClickMap       map[string]bool

	// 2 positions can be selected
	selectedCells [2]vec2i
}

func (grid *TileGrid) addSelection(coord vec2i) {
	coord.valid = true
	if !grid.selectedCells[0].valid {
		grid.selectedCells[0] = coord
	} else if !grid.selectedCells[1].valid {
		grid.selectedCells[1] = coord
	} else {
		// We could "rotate" the selections around, but in the context of the
		// game it doesn't really make sense. We'll just ignore it.
	}
}

func (grid *TileGrid) clearSelection() {
	grid.selectedCells[0] = vec2i{}
	grid.selectedCells[1] = vec2i{}
}

func drawGridTree(g *Game, tree *GridTree, screen *ebiten.Image, offsetY, offsetX int) {
	// Draw current node
	if tree.grid.SizeX == 0 {
		fmt.Println("GRID IS NIL", tree.grid)
		return
	}

	// Handle selected grid
    if g.selected != nil {
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
		faux := createGrid(offsetX, offsetY+tree.generation*120, tree.grid.SizeX, tree.grid.SizeY, 110, 110, tree.grid.Color)
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
        ClickMap: make(map[string]bool),
	}
	grid.Tiles = make([][]*Tile, grid.SizeX)
	for i := range grid.Tiles {
		grid.Tiles[i] = make([]*Tile, grid.SizeY)
	}

	for j := 0; j < grid.SizeY; j++ {
		for i := 0; i < grid.SizeX; i++ {
			if i%2 == 0 {
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

func (grid *TileGrid) Clone() TileGrid {
    ng := TileGrid{
        X: grid.X,
        Y: grid.Y,
        SizeX: grid.SizeX,
        SizeY: grid.SizeY,
        BoundsX: grid.BoundsX,
        BoundsY: grid.BoundsY,
        Color: grid.Color,
        IsSelectedGrid: grid.IsSelectedGrid,
        ClickMap: grid.ClickMap,
        selectedCells: grid.selectedCells,
    }
	ng.Tiles = make([][]*Tile, ng.SizeX)
	for i := range ng.Tiles {
		ng.Tiles[i] = make([]*Tile, ng.SizeY)
	}
    for j := 0; j < len(grid.Tiles); j++ {
        for i := 0; i < len(grid.Tiles[j]); i++ {
            x := *grid.Tiles[j][i]
            ng.Tiles[j][i] = &x
        }
    }

    return ng
}

func (grid *TileGrid) Update(g *Game) {
	if grid.IsSelectedGrid {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			grid.IsSelectedGrid = false
			g.selected = nil
			fmt.Println("KILLED")
		}
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            if !grid.ClickMap["clickTile"] {
                for j := 0; j < grid.SizeY; j++ {
                    for i := 0; i < grid.SizeX; i++ {
                        X, Y := tileScreenPos(grid, i, j)
                        r := tileRadius(grid)
                        mx, my := ebiten.CursorPosition()
                        if mx <= X+r && mx >= X-r && my <= Y+r && my >= Y-r {
                            grid.IsSelectedGrid = false
                            ng := grid.Clone()
                            ng.Tiles[j][i].Selected = true
                            gitCommitGrid(g, ng, false)
                            g.selected = &ng
                            g.selected.IsSelectedGrid = true
                            grid.ClickMap["clickTile"] = true
                            g.logger.AddMessage("you$ ", "git commit -m 'select a piece'", true)
                            g.logger.AddMessage("", "[main d34db33f] select a piece", true)
                            g.logger.AddMessage("", "1 files changed, 1 insertions(+), 0 deletions(-)", true)
                        }
                    }
                }
            }
        } else {
            grid.ClickMap["clickTile"] = false
        }

		if grid.selectedCells[1].valid {
			// User wants to make a move!
			pos1 := grid.selectedCells[0]
			pos2 := grid.selectedCells[1]
			source := grid.Tiles[pos1.x][pos1.y]
			target := grid.Tiles[pos2.x][pos2.y]

			if !source.occupant.Present {
				// No one's here, they can just move.
				target.occupant = source.occupant
				source.occupant = Unit{Present: false}
			} else {
				source.occupant.attackEnemy(&target.occupant)
				if target.occupant.HP <= 0 {
					// Attacker won, move into the cell
					target.occupant = source.occupant
					source.occupant = Unit{Present: false}
				} else if source.occupant.HP <= 0 {
					// Attacker lost, delete them
					source.occupant = Unit{Present: false}
				}
			}

			// The move is over, selection vanishes no matter what
			grid.clearSelection()
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
