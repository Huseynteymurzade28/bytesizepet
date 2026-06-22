package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
)

func colorToRGB(hex string) (colorful.Color, error) {
	return colorful.Hex(hex)
}

// mix linearly interpolates a single 0..1 channel toward target by t.
func mix(v, target, t float64) float64 {
	return v + (target-v)*t
}

// gradientBar renders a stat bar whose filled portion flows through a colour
// gradient (lo → hi) cell by cell, with a soft trailing glow character. This is
// the headline visual upgrade over the old single-colour bars.
func gradientBar(value, width int, lo, hi lipgloss.Color) string {
	if width <= 0 {
		return ""
	}
	loC, err1 := colorful.Hex(string(lo))
	hiC, err2 := colorful.Hex(string(hi))
	if err1 != nil || err2 != nil {
		loC, hiC = colorful.Color{R: 1, G: 1, B: 1}, colorful.Color{R: 1, G: 1, B: 1}
	}

	filled := value * width / 100
	if filled > width {
		filled = width
	}

	var b strings.Builder
	for i := 0; i < filled; i++ {
		t := 0.0
		if width > 1 {
			t = float64(i) / float64(width-1)
		}
		c := loC.BlendLab(hiC, t).Clamped()
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Render("█"))
	}
	// leading glow cell on the boundary
	if filled < width {
		t := 0.0
		if width > 1 {
			t = float64(filled) / float64(width-1)
		}
		c := loC.BlendLab(hiC, t).Clamped()
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Render("▓"))
		b.WriteString(barEmptyStyle.Render(strings.Repeat("░", width-filled-1)))
	}
	return b.String()
}
