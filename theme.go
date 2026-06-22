package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// dayPhase is the time-of-day bucket derived from the real system clock.
type dayPhase int

const (
	phaseDawn  dayPhase = iota // 5–8   soft pink/gold
	phaseDay                   // 8–17  bright sky
	phaseDusk                  // 17–20 orange/violet
	phaseNight                 // 20–5  deep indigo + stars
)

// theme bundles the colours and decorative sky used for one phase of the day.
type theme struct {
	name      string
	icon      string
	border    lipgloss.Color // frame border
	accent    lipgloss.Color // titles / highlights
	bg        lipgloss.Color // whole-stage atmosphere tint (shifts with the day)
	sky       lipgloss.Color // the strip of sky behind the pet
	skyGlyphs []rune         // scattered through the sky (stars, clouds, sun…)
	skyDense  int            // 1 in N sky cells gets a glyph
	ground    lipgloss.Color // habitat / floor colour
}

func currentPhase(now time.Time) dayPhase {
	switch h := now.Hour(); {
	case h >= 5 && h < 8:
		return phaseDawn
	case h >= 8 && h < 17:
		return phaseDay
	case h >= 17 && h < 20:
		return phaseDusk
	default:
		return phaseNight
	}
}

func themeFor(p dayPhase) theme {
	switch p {
	case phaseDawn:
		return theme{
			name: "Dawn", icon: "🌅",
			border: "#f3c0d0", accent: "#f7b8a0", bg: "#3a2e30", sky: "#f6d6c2",
			skyGlyphs: []rune{'☁', '·', '*'}, skyDense: 11, ground: "#cdb48e",
		}
	case phaseDay:
		return theme{
			name: "Day", icon: "☀️",
			border: "#bfe0f2", accent: "#f5d070", bg: "#26323a", sky: "#bfe6f5",
			skyGlyphs: []rune{'☁', '☁', '·'}, skyDense: 13, ground: "#a7c9a0",
		}
	case phaseDusk:
		return theme{
			name: "Dusk", icon: "🌇",
			border: "#d8a8c8", accent: "#f0a060", bg: "#332838", sky: "#caa0c0",
			skyGlyphs: []rune{'·', '*', '☁'}, skyDense: 9, ground: "#9a86a0",
		}
	default: // night
		return theme{
			name: "Night", icon: "🌙",
			border: "#8c8cc0", accent: "#c4b5e0", bg: "#1e1e36", sky: "#9a9ad0",
			skyGlyphs: []rune{'·', '*', '✦', '·', '*'}, skyDense: 6, ground: "#5a5680",
		}
	}
}

// starfield deterministically scatters sky glyphs across a width so the pattern
// is stable between frames (no flickering) but unique per width.
func (t theme) starfield(w int) []cell {
	cells := make([]cell, w)
	if len(t.skyGlyphs) == 0 || t.skyDense <= 0 {
		return cells
	}
	// simple LCG seeded by width keeps this self-contained and stable
	seed := uint32(w*2654435761 + 12345)
	next := func() uint32 { seed = seed*1664525 + 1013904223; return seed }
	for x := 0; x < w; x++ {
		if int(next()%uint32(t.skyDense)) == 0 {
			g := t.skyGlyphs[int(next())%len(t.skyGlyphs)]
			cells[x] = cell{g, lighten(t.sky)}
		}
	}
	return cells
}

// lighten nudges a sky colour brighter for the glyphs sprinkled on it.
func lighten(c lipgloss.Color) lipgloss.Color {
	col, err := colorToRGB(string(c))
	if err != nil {
		return c
	}
	col.R = mix(col.R, 1, 0.35)
	col.G = mix(col.G, 1, 0.35)
	col.B = mix(col.B, 1, 0.35)
	return lipgloss.Color(col.Hex())
}
