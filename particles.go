package main

import (
	"math/rand"

	"github.com/charmbracelet/lipgloss"
)

// particle is one floating glyph on the stage. Positions are floats so motion
// looks smooth even though the canvas is a coarse character grid.
type particle struct {
	x, y   float64
	vx, vy float64
	glyph  rune
	color  lipgloss.Color
	life   int // frames remaining; <=0 is culled
	wobble float64
}

// stepParticles advances every particle and drops the dead ones.
func stepParticles(ps []particle, w, h int) []particle {
	out := ps[:0]
	for _, p := range ps {
		p.x += p.vx
		p.y += p.vy
		p.wobble += 0.4
		p.life--
		if p.life <= 0 || p.y < -1 || p.y > float64(h)+1 || p.x < -2 || p.x > float64(w)+2 {
			continue
		}
		out = append(out, p)
	}
	return out
}

// spawnWeather seeds rain/snow falling from the top of the stage.
func spawnWeather(ps []particle, kind weatherKind, w int) []particle {
	switch kind {
	case weatherRain:
		if rand.Intn(2) == 0 {
			ps = append(ps, particle{
				x: float64(rand.Intn(w)), y: 0,
				vy: 1.1, vx: -0.15,
				glyph: '╱', color: "#8fb6e0", life: 14,
			})
		}
	case weatherSnow:
		if rand.Intn(3) == 0 {
			ps = append(ps, particle{
				x: float64(rand.Intn(w)), y: 0,
				vy: 0.35, vx: 0,
				glyph: '*', color: "#eef4ff", life: 30, wobble: rand.Float64() * 6,
			})
		}
	case weatherThunder:
		if rand.Intn(2) == 0 {
			ps = append(ps, particle{
				x: float64(rand.Intn(w)), y: 0,
				vy: 1.4, vx: -0.2,
				glyph: '╱', color: "#aab0d0", life: 12,
			})
		}
	}
	return ps
}

// spawnMood emits the little emotional flourishes around the pet: hearts when
// happy, music notes while playing, sweat when starving.
func spawnMood(ps []particle, m model, centerX, topY int) []particle {
	p := m.pet
	switch {
	case m.actionAnim == 2: // playing → hearts burst
		if rand.Intn(2) == 0 {
			ps = append(ps, particle{
				x: float64(centerX - 2 + rand.Intn(5)), y: float64(topY + 1),
				vy: -0.5, vx: (rand.Float64() - 0.5) * 0.4,
				glyph: '♥', color: heartColor(), life: 12,
			})
		}
	case p.sleeping:
		// handled by the floating zZz in the renderer
	case p.happiness >= 80 && !p.sleeping:
		if rand.Intn(4) == 0 {
			ps = append(ps, particle{
				x: float64(centerX - 3 + rand.Intn(7)), y: float64(topY),
				vy: -0.4, vx: (rand.Float64() - 0.5) * 0.3,
				glyph: pick('♥', '♡'), color: heartColor(), life: 14,
			})
		}
	case p.hunger >= 75:
		if rand.Intn(5) == 0 {
			ps = append(ps, particle{
				x: float64(centerX + 4), y: float64(topY + 1),
				vy: 0.5, vx: 0.1,
				glyph: '✦', color: "#86c8e8", life: 8,
			})
		}
	}
	return ps
}

func pick(opts ...rune) rune { return opts[rand.Intn(len(opts))] }

func heartColor() lipgloss.Color {
	cols := []lipgloss.Color{"#f08fb4", "#f7a8c8", "#ff6f91", "#f5c0d8"}
	return cols[rand.Intn(len(cols))]
}
