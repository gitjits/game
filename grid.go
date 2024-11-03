package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed hexagon.png
var hexagonPng []byte
var hexagonImg *ebiten.Image

//go:embed lbj.png
var LBJPNG []byte
var LBJ_Img *ebiten.Image

func loadEmbeddedImage() (err error) {
	hi, _, err := image.Decode(bytes.NewReader(hexagonPng))
	lbj_img, _, err := image.Decode(bytes.NewReader(LBJPNG))
	if err != nil {
		return err
	}
	hexagonImg = ebiten.NewImageFromImage(hi)
	LBJ_Img = ebiten.NewImageFromImage(lbj_img)
	return nil
}

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

func (grid *TileGrid) posSelected(coord vec2i) bool {
	cells := grid.selectedCells
	if cells[0].valid && cells[0].x == coord.x && cells[0].y == coord.y {
		return true
	}
	if cells[1].valid && cells[1].x == coord.x && cells[1].y == coord.y {
		return true
	}
	return false
}

func (grid *TileGrid) clearSelection() {
	grid.selectedCells[0] = vec2i{}
	grid.selectedCells[1] = vec2i{}
}

func (grid *TileGrid) applyMove() {
	if !grid.selectedCells[1].valid {
		return
	}

	// User wants to make a move!
	pos1 := grid.selectedCells[0]
	pos2 := grid.selectedCells[1]
	source := grid.Tiles[pos1.x][pos1.y]
	target := grid.Tiles[pos2.x][pos2.y]

	if !source.occupant.Present {
		// No one's here, they can just move.
		*target = *source
		source.occupant = Unit{Present: false}
		source.Color = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
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

func drawGridTree(g *Game, tree *GridTree, screen *ebiten.Image, offsetY, offsetX int) {
	// Draw current node
	if tree.grid.SizeX == 0 {
		fmt.Println("GRID IS NIL", tree.grid)
		return
	}

	// Handle selected grid
	if g.selected != nil {
		g.selected.X = screenWidth/2 - 150
		g.selected.Y = screenHeight/2 - 150
		g.selected.BoundsX = 310
		g.selected.BoundsY = 340
		g.selected.Update(g)
		r := tileRadius(g.selected)
		vector.DrawFilledRect(screen, float32(g.selected.X-r/2), float32(g.selected.Y), float32(g.selected.BoundsX+r), float32(g.selected.BoundsY+r), color.RGBA{0, 0, 0, 100}, false)
		if g.selected != nil {
			drawGrid(*g.selected, screen)
		}
	}
	if tree.grid.IsSelectedGrid {
		// Draw main selected grid in center

		// Draw small version in tree
		faux := createGrid(offsetX, offsetY+tree.generation*125, tree.grid.SizeX, tree.grid.SizeY, 115, 123, tree.grid.Color)
		faux.Tiles = tree.grid.Tiles
		drawGrid(faux, screen)
	} else {
		tree.grid.X = offsetX
		tree.grid.Y = offsetY + tree.generation*125
		tree.grid.BoundsX = 115
		tree.grid.BoundsY = 123
	}

	tree.grid.Update(g)
	drawGrid(tree.grid, screen)

	// Continue main branch
	if tree.next != nil && tree.next.grid.SizeX != 0 {
		drawGridTree(g, tree.next, screen, offsetY, offsetX+135)
	}
}

func drawGrid(grid TileGrid, screen *ebiten.Image) {
	r := tileRadius(&grid)
	for j := 0; j < len(grid.Tiles); j++ {
		for i := 0; i < len(grid.Tiles[j]); i++ {
			tile := grid.Tiles[j][i]
			Xpos, Ypos := tileScreenPos(&grid, i, j)
			op := &colorm.DrawImageOptions{}
			scale := float64(float64(r) * 2.0 / 256.0)
			op.GeoM.Scale(1.25*scale, scale*1.05)
			op.GeoM.Translate(float64(Xpos-r), float64(Ypos-r))
			var cm colorm.ColorM
			r := float64(tile.Color.R) / 0xff
			g := float64(tile.Color.G) / 0xff
			b := float64(tile.Color.B) / 0xff
			a := 1.0
			if grid.posSelected(vec2i{x: j, y: i}) {
				a = 0.25
			}
			cm.Scale(r, g, b, a)
			colorm.DrawImage(screen, hexagonImg, cm, op)

			unitOptions := &ebiten.DrawImageOptions{}
			unitOptions.GeoM = op.GeoM
			if tile.occupant.Name == UNIT_LBJ.Name {
				screen.DrawImage(LBJ_Img, unitOptions)
			}
		}
	}
	vector.StrokeRect(screen, float32(grid.X-r/2), float32(grid.Y), float32(grid.BoundsX+r), float32(grid.BoundsY+r), 1, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)
}

func createGrid(X int, Y int, SizeX int, SizeY int, BoundsX int, BoundsY int, defaultColor color.RGBA) TileGrid {
	grid := TileGrid{
		X:        X,
		Y:        Y,
		SizeX:    SizeX,
		SizeY:    SizeY,
		BoundsX:  BoundsX,
		BoundsY:  BoundsY,
		Color:    defaultColor,
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

	var X = grid.X + 2 + w*row + r
	var Y = grid.Y + 3 + int(float64(h)*float64(col)*1.1) + int(float64(r)*ymult)

	return X, Y
}

func (grid *TileGrid) Clone() TileGrid {
	ng := TileGrid{
		X:              grid.X,
		Y:              grid.Y,
		SizeX:          grid.SizeX,
		SizeY:          grid.SizeY,
		BoundsX:        grid.BoundsX,
		BoundsY:        grid.BoundsY,
		Color:          grid.Color,
		IsSelectedGrid: grid.IsSelectedGrid,
		ClickMap:       grid.ClickMap,
		selectedCells:  grid.selectedCells,
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
							grid.addSelection(vec2i{x: j, y: i})
							grid.ClickMap["clickTile"] = true
						}
					}
				}
			}
		} else {
			grid.ClickMap["clickTile"] = false
		}

		grid.applyMove()
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
