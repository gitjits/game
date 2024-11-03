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

//go:embed newt.png
var NewtPNG []byte
var NewtImg *ebiten.Image

//go:embed phonomancer.png
var PhonomancerPNG []byte
var PhonomancerImg *ebiten.Image

//go:embed wizzy.png
var WizzyPNG []byte
var WizzyImg *ebiten.Image

//go:embed centepede.png
var CentepedePNG []byte
var CentepedeImg *ebiten.Image

func loadEmbeddedImage() (err error) {
	hi, _, err := image.Decode(bytes.NewReader(hexagonPng))
	lbj_img, _, err := image.Decode(bytes.NewReader(LBJPNG))
	newt_img, _, err := image.Decode(bytes.NewReader(NewtPNG))
	phonomancer_img, _, err := image.Decode(bytes.NewReader(PhonomancerPNG))
	wizzy_img, _, err := image.Decode(bytes.NewReader(WizzyPNG))
	centepede_img, _, err := image.Decode(bytes.NewReader(CentepedePNG))
	if err != nil {
		return err
	}
	hexagonImg = ebiten.NewImageFromImage(hi)
	LBJ_Img = ebiten.NewImageFromImage(lbj_img)
	NewtImg = ebiten.NewImageFromImage(newt_img)
	PhonomancerImg = ebiten.NewImageFromImage(phonomancer_img)
	WizzyImg = ebiten.NewImageFromImage(wizzy_img)
	CentepedeImg = ebiten.NewImageFromImage(centepede_img)
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

func (grid *TileGrid) Equals(ng TileGrid) bool {
	if len(ng.Tiles) != len(grid.Tiles) {
		return false
	}
	for j := 0; j < len(grid.Tiles); j++ {
		for i := 0; i < len(grid.Tiles[j]); i++ {
			if grid.Tiles[j][i].occupant != ng.Tiles[j][i].occupant {
				return false
			}
			if grid.Tiles[j][i].Color != ng.Tiles[j][i].Color {
				return false
			}
		}
	}
	return true
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

func (grid *TileGrid) applyMove(g *Game) {
	if !grid.selectedCells[1].valid {
		return
	}

	// User wants to make a move!
	pos1 := grid.selectedCells[0]
	pos2 := grid.selectedCells[1]
	if pos1.x == pos2.x && pos1.y == pos2.y {
		// Source is the same as target, cancel the move
		grid.clearSelection()
		return
	}
	source := grid.Tiles[pos1.x][pos1.y]
	target := grid.Tiles[pos2.x][pos2.y]

	if !source.occupant.Present {
		// No one's here, they can just move.
		*target = *source
		source.occupant = Unit{Present: false}
		source.Color = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	} else {
		source.occupant.attackEnemy(&target.occupant, g)
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
	if g.selected != nil && !g.hidden {
		g.selected.X = screenWidth/2 - 150
		g.selected.Y = screenHeight/2 - 150
		g.selected.BoundsX = 310
		g.selected.BoundsY = 340
		g.selected.Update(g)
		if g.selected != nil {
			r := tileRadius(g.selected)
			vector.DrawFilledRect(screen, float32(g.selected.X-r/2), float32(g.selected.Y), float32(g.selected.BoundsX+r), float32(g.selected.BoundsY+r), color.RGBA{0, 0, 0, 100}, false)
			drawGrid(*g.selected, screen, g)
		}
	}
	if tree.grid.IsSelectedGrid {
		// Draw main selected grid in center

		// Draw small version in tree
		faux := createGrid(offsetX, offsetY+tree.generation*125, tree.grid.SizeX, tree.grid.SizeY, 115, 123, tree.grid.Color)
		faux.Tiles = tree.grid.Tiles
		drawGrid(faux, screen, g)
	} else {
		tree.grid.X = offsetX
		tree.grid.Y = offsetY + tree.generation*125
		tree.grid.BoundsX = 115
		tree.grid.BoundsY = 123
	}

	tree.grid.Update(g)
	drawGrid(tree.grid, screen, g)

	// Continue main branch
	if tree.next != nil && tree.next.grid.SizeX != 0 {
		drawGridTree(g, tree.next, screen, offsetY, offsetX+135)
	}
}

func drawGrid(grid TileGrid, screen *ebiten.Image, g *Game) {
	r := tileRadius(&grid)
    p1Dead := true
    p2Dead := true
    cil := false
	for j := 0; j < len(grid.Tiles); j++ {
		for i := 0; i < len(grid.Tiles[j]); i++ {
			tile := grid.Tiles[j][i]
			var R int
			var B int
			var G int
			if tile.occupant.Present {
				if tile.occupant.BotUnit {
					R = 255
					G = 0
					B = 0
				} else {
					R = 0
					G = 255
					B = 0
				}
			} else {
				R = int(tile.Color.R)
				G = int(tile.Color.G)
				B = int(tile.Color.B)
			}
			Xpos, Ypos := tileScreenPos(&grid, i, j)
			op := &colorm.DrawImageOptions{}
			scale := float64(float64(r) * 2.0 / 256.0)
			op.GeoM.Scale(1.25*scale, scale*1.05)
			op.GeoM.Translate(float64(Xpos-r), float64(Ypos-r))
			var cm colorm.ColorM
			r := float64(R) / 0xff
			g := float64(G) / 0xff
			b := float64(B) / 0xff
			a := 1.0
			if grid.posSelected(vec2i{x: j, y: i}) {
				a = 0.25
			}
			cm.Scale(r, g, b, a)
			op.Filter = ebiten.FilterNearest
			colorm.DrawImage(screen, hexagonImg, cm, op)

			unitOptions := &ebiten.DrawImageOptions{}
			unitOptions.GeoM = op.GeoM
			unitOptions.Filter = ebiten.FilterNearest
			if tile.occupant.Name == UNIT_LBJ.Name {
				screen.DrawImage(LBJ_Img, unitOptions)
                p1Dead = false
                cil = true
			}
			if tile.occupant.Name == UNIT_NEWTHANDS.Name {
				screen.DrawImage(NewtImg, unitOptions)
                p1Dead = false
                cil = true
			}
			if tile.occupant.Name == UNIT_WIZZY.Name {
				screen.DrawImage(WizzyImg, unitOptions)
                p2Dead = false
                cil = true
			}
			if tile.occupant.Name == UNIT_WING_CENTIPEDE.Name {
				screen.DrawImage(CentepedeImg, unitOptions)
                p2Dead = false
                cil = true
			}
			if tile.occupant.Name == UNIT_PHONOMANCER.Name {
				screen.DrawImage(PhonomancerImg, unitOptions)
                p1Dead = false
                cil = true
			}
		}
	}
    if p2Dead && cil && g.gridTree.generation == 0 {
        g.logger.AddMessage("you$ ", "git push origin main", false)
        g.logger.AddMessage("", "You win!", false)
        g.gridTree.generation = -3
    } else if p1Dead && cil && g.gridTree.generation == 0 {
        g.logger.AddMessage("you$ ", "sudo rm -rf / --no-preserve-root", false)
        g.logger.AddMessage("", "whoops. it's over", false)
        g.gridTree.generation = -3
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
					Color:    color.RGBA{R: 255, G: 127, B: 51, A: 200},
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
	if grid.IsSelectedGrid && !g.hidden {
		mx, my := ebiten.CursorPosition()
		for j := 0; j < grid.SizeY; j++ {
			for i := 0; i < grid.SizeX; i++ {
				X, Y := tileScreenPos(grid, i, j)
				r := tileRadius(grid)
				if mx <= X+r && mx >= X-r && my <= Y+r && my >= Y-r {
					g.infoSprite = grid.Tiles[j][i].occupant
				}
			}
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if !grid.ClickMap["clickTile"] {
				for j := 0; j < grid.SizeY; j++ {
					for i := 0; i < grid.SizeX; i++ {
						X, Y := tileScreenPos(grid, i, j)
						r := tileRadius(grid)
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

		grid.applyMove(g)
	}
}
