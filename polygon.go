package main

import (
	"math"
	_ "embed"
    "image"
    "image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func genVertices(num int, centerX float32, centerY float32, r float64, hc color.RGBA) []ebiten.Vertex {
	vs := []ebiten.Vertex{}
	for i := 0; i < num; i++ {
		rate := float64(i) / float64(num)
		vs = append(vs, ebiten.Vertex{
			DstX:   float32(r*math.Cos(2*math.Pi*rate)) + centerX,
			DstY:   float32(r*math.Sin(2*math.Pi*rate)) + centerY,
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(hc.R),
			ColorG: float32(hc.G),
			ColorB: float32(hc.B),
			ColorA: 1,
		})
	}

	vs = append(vs, ebiten.Vertex{
		DstX:   centerX,
		DstY:   centerY,
		SrcX:   0,
		SrcY:   0,
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	})

	return vs
}

func drawPolygon(n int, x int, y int, r int, hc color.RGBA, screen *ebiten.Image) {
    op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	indices := []uint16{}
	for i := 0; i < n; i++ {
		indices = append(indices, uint16(i), uint16(i+1)%uint16(n), uint16(n))
	}
    whiteImage := ebiten.NewImage(3, 3)
    whiteImage.Fill(hc)
	screen.DrawTriangles(genVertices(n, float32(x), float32(y), float64(r), hc), indices, whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
}
