package main

import (
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

type LogMessage struct {
    Text      string
    Prompt string
    Done bool
    Timestamp time.Time
}

type LogWindow struct {
    messages    []LogMessage
    maxMessages int
    lineHeight  float64
    font        font.Face
    fadeHeight  float64
    status      int
    prompt string
    height      float64
}

func NewLogWindow() *LogWindow {
    return &LogWindow{
        messages:    make([]LogMessage, 0),
        maxMessages: 100,
        lineHeight: 20,
        font:       inconsolata.Regular8x16,
        fadeHeight: 50,
        status: 0,
        height: 250,
    }
}

func (l *LogWindow) AddMessage(prompt, msg string, instant bool) {
    l.messages = append(l.messages, LogMessage{
        Text:      msg,
        Prompt: prompt,
        Done: false,
        Timestamp: time.Now(),
    })
    if len(l.messages) > l.maxMessages {
        l.messages = l.messages[1:]
    }
    if instant {
        l.status = 0
    }
}

func (l *LogWindow) Draw(screen *ebiten.Image) {
    height := l.height
    width := float64(screenWidth/2)
    sx := float64(0)
    sy := float64(screenHeight)-l.height
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

    l.status += 1
    if l.status == 1000 {
        l.height = 150
    }
    ch := 0
    for i, msg := range l.messages[startIdx:] {
        yPos := sy + float64(i)*l.lineHeight
        alpha := uint8(255)
        if yPos < sy+l.fadeHeight {
            alpha = uint8((yPos - sy) / l.fadeHeight * 255)
        }
        col := color.RGBA{0, 255, 50, alpha}
        cum := msg.Prompt
        for _, char := range strings.Split(msg.Text, "") {
            if !msg.Done {
                ch += 1
            }
            cum += char
            if ch > l.status && !msg.Done {
                break
            }
            text.Draw(screen, cum, l.font, int(sx)+15, int(yPos)+15, col)
        }
        if i == len(l.messages[startIdx:]) - 1 && len(cum) >= len(msg.Prompt) + len(msg.Text) {
            for xx, _ := range l.messages {
                l.messages[xx].Done = true
            }
        }
    }
}
