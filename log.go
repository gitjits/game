package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

type LogMessage struct {
    Text      string
    Timestamp time.Time
}

type LogWindow struct {
    messages    []LogMessage
    maxMessages int
    lineHeight  float64
    font        font.Face
    fadeHeight  float64
}

func NewLogWindow() *LogWindow {
    return &LogWindow{
        messages:    make([]LogMessage, 0),
        maxMessages: 100,
        lineHeight: 20,
        font:       inconsolata.Regular8x16,
        fadeHeight: 50,
    }
}

func (l *LogWindow) AddMessage(msg string) {
    l.messages = append(l.messages, LogMessage{
        Text:      msg,
        Timestamp: time.Now(),
    })
    if len(l.messages) > l.maxMessages {
        l.messages = l.messages[1:]
    }
}

func (l *LogWindow) Draw(screen *ebiten.Image) {
    height := float64(screenHeight/5)
    width := float64(screenHeight/2)
    sx := float64(0)
    sy := float64(screenHeight*4/5)
    for y := float64(0); y < height; y++ {
        alpha := uint8(255)
        if y < l.fadeHeight {
            alpha = uint8(float64(y) / l.fadeHeight * 255)
        }
        col := color.RGBA{0, 0, 0, alpha}
        ebitenutil.DrawRect(screen, sx, sy+y, width, 1, col)
    }

    visibleLines := int(height / l.lineHeight)
    startIdx := len(l.messages)
    if startIdx > visibleLines {
        startIdx = len(l.messages) - visibleLines
    }

    for i, msg := range l.messages[startIdx:] {
        yPos := sy + float64(i)*l.lineHeight
        alpha := uint8(255)
        if yPos < sy+l.fadeHeight {
            alpha = uint8((yPos - sy) / l.fadeHeight * 255)
        }
        col := color.RGBA{0, 255, 50, alpha}
        text.Draw(screen, msg.Text, l.font, int(sx)+15, int(yPos)+15, col)
    }
}
