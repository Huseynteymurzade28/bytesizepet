package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// cell is a single character on the stage. An empty rune means "transparent"
// and lets whatever was drawn earlier (or the background) show through.
type cell struct {
	r     rune
	color lipgloss.Color
}

// canvas is a small fixed-size grid we composite the scene onto: sky, then the
// habitat, then the pet, then weather + mood particles on top. Rendering walks
// each row and coalesces runs of identical colour so we emit as little ANSI as
// possible per frame.
type canvas struct {
	w, h  int
	cells []cell
}

func newCanvas(w, h int) *canvas {
	return &canvas{w: w, h: h, cells: make([]cell, w*h)}
}

func (c *canvas) inBounds(x, y int) bool {
	return x >= 0 && x < c.w && y >= 0 && y < c.h
}

func (c *canvas) set(x, y int, r rune, color lipgloss.Color) {
	if !c.inBounds(x, y) || r == ' ' || r == 0 {
		return
	}
	c.cells[y*c.w+x] = cell{r, color}
}

// drawString stamps a single line of text starting at (x,y). Spaces are treated
// as transparent so sprites don't punch holes in the background.
func (c *canvas) drawString(x, y int, s string, color lipgloss.Color) {
	col := x
	for _, r := range s {
		c.set(col, y, r, color)
		col++
	}
}

// drawLinesCentered stamps a block of text horizontally centred, top-aligned at y.
func (c *canvas) drawLinesCentered(y int, lines []string, color lipgloss.Color) {
	for i, ln := range lines {
		x := (c.w - lipgloss.Width(ln)) / 2
		c.drawString(x, y+i, ln, color)
	}
}

// render flattens the grid into a styled string.
func (c *canvas) render(bg lipgloss.Color) string {
	var b strings.Builder
	for y := 0; y < c.h; y++ {
		var runStr strings.Builder
		runColor := bg
		flush := func() {
			if runStr.Len() == 0 {
				return
			}
			b.WriteString(lipgloss.NewStyle().Foreground(runColor).Render(runStr.String()))
			runStr.Reset()
		}
		for x := 0; x < c.w; x++ {
			cl := c.cells[y*c.w+x]
			r := cl.r
			color := cl.color
			if r == 0 {
				r = ' '
				color = bg
			}
			if color != runColor {
				flush()
				runColor = color
			}
			runStr.WriteRune(r)
		}
		flush()
		if y < c.h-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
