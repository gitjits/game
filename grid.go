package main

import (
	//"fmt"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	width  int
	height int
	x      int
	y      int
    color  color.RGBA
    radius int
}

type TileGrid struct {
	tiles [][]*Tile
	x     int
	y     int
    sizeX int
    sizeY int
    boundsX int
    boundsY int
}

func drawGrid(grid TileGrid, screen *ebiten.Image) {
	for j := 0; j < len(grid.tiles); j++ {
		for i := 0; i < len(grid.tiles[j]); i++ {
			tile := grid.tiles[j][i]
			drawPolygon(6, tile.x+grid.x, tile.y+grid.y, tile.radius, tile.color, screen)
		}
	}
}

func createGrid(x int, y int, sizeX int, sizeY int, boundsX int, boundsY int) TileGrid {
    grid := TileGrid{
        x: x,
        y: y,
        sizeX: sizeX,
        sizeY: sizeY,
        boundsX: boundsX,
        boundsY: boundsY,
    }
	grid.tiles = make([][]*Tile, grid.sizeX)
	for i := range grid.tiles {
		grid.tiles[i] = make([]*Tile, grid.sizeY)
	}

	for j := 0; j < grid.sizeY; j++ {
		for i := 0; i < grid.sizeX; i++ {
            grid.tiles[j][i] = &Tile{
                color: color.RGBA{R: 255, B: 255, G: 255, A: 1},
            }
        }
    }
    grid.Update()

    return grid
}

func (grid *TileGrid) Update() {
	for j := 0; j < grid.sizeY; j++ {
		for i := 0; i < grid.sizeX; i++ {
            r := grid.boundsX / grid.sizeX / 2
			w, h := r*2, r*2
			ymult := 1.0
			if i%2 != 0 {
				ymult = 2
			}
			x, y := w*i+r, h*j+int(math.Floor(float64(r)*ymult))
            grid.tiles[j][i].width = w
            grid.tiles[j][i].height = h
            grid.tiles[j][i].x = x
            grid.tiles[j][i].y = y
            grid.tiles[j][i].radius = r
            fmt.Println(grid.tiles[j][i])
		}
	}
}
