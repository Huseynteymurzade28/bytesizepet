package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Stage geometry — the little diorama the pet lives in.
const (
	stageW  = 40
	stageH  = 15
	petTopY = 3 // first row of the sprite block
)

// ── Root view ─────────────────────────────────────────────────────────────────

func (m model) View() string {
	th := m.theme()

	if m.quitting {
		msg := "Come back soon~ 🐾"
		if m.pet.name != "" {
			msg = m.pet.name + " will miss you~ 🐾"
		}
		farewell := lipgloss.NewStyle().Foreground(th.accent).
			Render("\n  Goodbye! " + msg + "\n\n")
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, farewell)
	}

	var content string
	switch m.screen {
	case screenCharSelect:
		content = viewCharSelect(m, th)
	case screenMain:
		content = viewMain(m, th)
	case screenMinigame:
		content = viewMinigame(m, th)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func frame(th theme) lipgloss.Style {
	return borderStyle.BorderForeground(th.border)
}

// ── Character select ──────────────────────────────────────────────────────────

func viewCharSelect(m model, th theme) string {
	title := titleStyle.Foreground(th.accent).Render("✨  Choose Your Pet  ✨")
	sub := helpStyle.Render("  ← / → / ↑ / ↓  to select,  Enter to confirm")

	boxes := make([]string, len(AllChars))
	for i, ch := range AllChars {
		art := lipgloss.NewStyle().Foreground(ch.Body).Render(
			strings.Join(ch.spriteLines(moodHappy), "\n"))
		nameStr := ch.Name
		if i == m.selectedChar {
			nameStr = lipgloss.NewStyle().Bold(true).Foreground(colorRose).Render("✦ " + nameStr + " ✦")
		}
		label := nameStr + "\n" + ch.Emoji + " " + ch.Species
		inner := art + "\n" + label

		if i == m.selectedChar {
			boxes[i] = charBoxSelectedStyle.Render(inner)
		} else {
			boxes[i] = charBoxStyle.Render(inner)
		}
	}

	var rows []string
	for i := 0; i < len(boxes); i += charSelectCols {
		end := i + charSelectCols
		if end > len(boxes) {
			end = len(boxes)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, boxes[i:end]...))
	}

	content := strings.Join(append([]string{title, sub, ""}, rows...), "\n")
	return frame(th).Render(content)
}

// ── Main screen ───────────────────────────────────────────────────────────────

func viewMain(m model, th theme) string {
	p := m.pet

	title := titleStyle.Foreground(th.accent).
		Render(fmt.Sprintf("✨  %s  %s  ✨", p.name, AllChars[p.charIdx].Emoji))

	stage := petAreaStyle.Render(renderStage(m, th))

	worldRow := renderWorldBar(m, th)

	stats := strings.Join([]string{
		renderStat("Hunger   :", p.hunger, "#8fd694", "#e06868"),
		renderStat("Happiness:", p.happiness, "#e07868", "#7ad6a0"),
		renderStat("Energy   :", p.energy, "#e07868", "#78c8e8"),
	}, "\n")

	center := lipgloss.NewStyle().Width(56).Align(lipgloss.Center)

	statusBox := center.Render(statusBoxStyle.Render(statusIcon(p) + " " + p.statusMsg))

	divider := center.Render(dividerStyle.Render(strings.Repeat("─", 44)))
	controls1 := center.Render(renderControls(
		[]string{"f", "p", "s", "n", "q"},
		[]string{"Feed", "Play", "Sleep", "New", "Quit"},
	))
	controls2 := center.Render(renderControls(
		[]string{"1", "2", "3", "4"},
		[]string{"Guess", "React", "RPS", "Math"},
	))
	statsBlock := center.Render(stats)

	content := strings.Join([]string{
		title,
		stage,
		worldRow, "",
		statsBlock, "",
		statusBox,
		divider,
		controls1,
		controls2,
	}, "\n")

	return frame(th).Render(content)
}

// renderWorldBar shows the time-of-day and live weather pulled from the API.
func renderWorldBar(m model, th theme) string {
	left := lipgloss.NewStyle().Foreground(th.accent).Render(th.icon + " " + th.name)
	dot := helpStyle.Render("  ·  ")

	if !m.weather.ok {
		right := helpStyle.Render("fetching weather…")
		return petAreaStyle.Render(left + dot + right)
	}
	w := m.weather
	right := lipgloss.NewStyle().Foreground(colorSkyBlue).
		Render(fmt.Sprintf("%s %s %d°C", w.kind.icon(), w.desc, w.tempC))
	area := ""
	if w.area != "" {
		area = dot + helpStyle.Render(w.area)
	}
	return petAreaStyle.Render(left + dot + right + area)
}

// ── Minigame screen ───────────────────────────────────────────────────────────

func viewMinigame(m model, th theme) string {
	title := titleStyle.Foreground(th.accent).
		Render(fmt.Sprintf("✨  %s  %s  ✨", m.pet.name, AllChars[m.pet.charIdx].Emoji))
	stage := petAreaStyle.Render(renderStage(m, th))

	var gameSection string
	switch {
	case m.mg.phase == mgGuessPlaying || m.mg.phase == mgGuessDone:
		gameSection = viewGuessGame(m)
	case m.mg.phase == mgRPSPlaying || m.mg.phase == mgRPSDone:
		gameSection = viewRPSGame(m)
	case m.mg.phase == mgMathPlaying || m.mg.phase == mgMathDone:
		gameSection = viewMathGame(m)
	default:
		gameSection = viewReactionGame(m)
	}

	help := helpStyle.Render("[ESC] Back to main")

	content := strings.Join([]string{
		title,
		stage, "",
		gameSection, "",
		help,
	}, "\n")

	return frame(th).Render(content)
}

func viewGuessGame(m model) string {
	header := mgTitleStyle.Render("🎮  NUMBER GUESS")
	if m.mg.phase == mgGuessDone {
		return strings.Join([]string{
			header, "",
			resultStyle.Render(m.mg.resultMsg), "",
			helpStyle.Render("[ENTER] Continue"),
		}, "\n")
	}
	hint := statusStyle.Render(m.mg.hint)
	keys := helpStyle.Render("[1-9] Make a guess")
	return strings.Join([]string{header, "", hint, keys}, "\n")
}

func viewReactionGame(m model) string {
	header := mgTitleStyle.Render("⚡  REACTION TIME")
	var body string
	switch m.mg.phase {
	case mgReactionWait:
		body = reactionWaitStyle.Render("Get ready... waiting for the signal...") +
			"\n\n" + helpStyle.Render("[SPACE] wait — don't press early!")
	case mgReactionReady:
		body = reactionReadyStyle.Render("  ★  NOW!  ★  ") +
			"\n\n" + helpStyle.Render("[SPACE] PRESS IT!")
	case mgReactionDone:
		body = resultStyle.Render(m.mg.resultMsg) +
			"\n\n" + helpStyle.Render("[ENTER] Continue")
	}
	return strings.Join([]string{header, "", body}, "\n")
}

func viewRPSGame(m model) string {
	header := mgRPSTitleStyle.Render("🪨  ROCK  PAPER  SCISSORS  ✂️")
	if m.mg.phase == mgRPSDone {
		lines := strings.SplitN(m.mg.resultMsg, "\n\n", 2)
		var body string
		if len(lines) == 2 {
			body = statusStyle.Render(lines[0]) + "\n\n" + resultStyle.Render(lines[1])
		} else {
			body = resultStyle.Render(m.mg.resultMsg)
		}
		return strings.Join([]string{header, "", body, "", helpStyle.Render("[ENTER] Continue")}, "\n")
	}
	prompt := statusStyle.Render("Choose your move!")
	choices := helpStyle.Render("[R] Rock 🪨   [P] Paper 📄   [S] Scissors ✂️")
	return strings.Join([]string{header, "", prompt, choices}, "\n")
}

func viewMathGame(m model) string {
	header := mgMathTitleStyle.Render("🧮  MATH QUIZ")
	if m.mg.phase == mgMathDone {
		return strings.Join([]string{
			header, "",
			resultStyle.Render(m.mg.resultMsg), "",
			helpStyle.Render("[ENTER] Continue"),
		}, "\n")
	}
	question := mathQuestionStyle.Render("  " + m.mg.mathExpr + "  ")
	opts := helpStyle.Render(fmt.Sprintf("[1] %-5d  [2] %-5d  [3] %d",
		m.mg.mathOptions[0], m.mg.mathOptions[1], m.mg.mathOptions[2]))
	return strings.Join([]string{header, "", question, "", opts}, "\n")
}

// ── The diorama ───────────────────────────────────────────────────────────────

// renderStage composites the whole scene: time-of-day sky, the ground, the pet
// sprite, sleep z's, the feeding fish, and every weather/mood particle.
func renderStage(m model, th theme) string {
	c := newCanvas(stageW, stageH)

	// 1. Sky — scattered stars/clouds on the top rows (stable per width).
	stars := th.starfield(stageW)
	for x, cl := range stars {
		if cl.r != 0 {
			c.set(x, 0, cl.r, cl.color)
		}
	}
	stars2 := th.starfield(stageW + 7) // different pattern for the 2nd row
	for x := 0; x < stageW; x++ {
		if cl := stars2[x]; cl.r != 0 {
			c.set(x, 1, cl.r, cl.color)
		}
	}

	// 2. Ground — a soft scrolling floor line with little tufts.
	groundY := stageH - 2
	groundGlyphs := []rune{'.', ',', '˙', '·'}
	for x := 0; x < stageW; x++ {
		g := groundGlyphs[(x+m.envFrame)%len(groundGlyphs)]
		c.set(x, groundY, g, th.ground)
	}
	for x := 0; x < stageW; x++ {
		c.set(x, groundY+1, '▁', th.ground)
	}

	// 3. The pet, centred and mood-tinted.
	mood := m.mood()
	ch := AllChars[m.pet.charIdx]
	sprite := ch.spriteLines(mood)
	tint := ch.moodColor(mood)
	// gentle idle bob: rise one row on alternate idle frames when content
	bob := 0
	if mood == moodHappy && m.animFrame == 1 {
		bob = -1
	}
	c.drawLinesCentered(petTopY+bob, sprite, tint)
	cx := stageW / 2

	// 4. Sleep z's drifting up from the head.
	if m.pet.sleeping {
		zs := []rune{'z', 'Z', 'z'}
		for i := 0; i <= m.sleepFrame && i < len(zs); i++ {
			c.set(cx+4+i, petTopY-i, zs[i%len(zs)], "#b8a8e8")
		}
	}

	// 5. Feeding fish swims in toward the mouth.
	if m.actionAnim == 1 {
		fishX := 4 + m.actionFrame*2
		fy := petTopY + 4
		if m.actionFrame >= 6 {
			c.drawString(cx-2, fy, "nom♪", colorSalmon)
		} else {
			c.drawString(fishX, fy, "<><", colorSalmon)
		}
	}

	// 6. Particles on top of everything.
	for _, p := range m.particles {
		dx := 0.0
		if p.glyph == '*' { // snow drifts side to side
			dx = math.Sin(p.wobble)
		}
		c.set(int(p.x+dx+0.5), int(p.y+0.5), p.glyph, p.color)
	}

	return c.render(th.bg)
}

// ── Control buttons ───────────────────────────────────────────────────────────

func renderControls(keys, labels []string) string {
	dot := helpStyle.Render(" · ")
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = keyStyle.Render("["+k+"]") + helpStyle.Render(" "+labels[i])
	}
	return strings.Join(parts, dot)
}

func statusIcon(p petModel) string {
	switch {
	case p.sleeping:
		return "💤"
	case p.hunger >= 75:
		return "🍽"
	case p.happiness <= 20:
		return "😢"
	case p.energy <= 20:
		return "😴"
	default:
		return "✨"
	}
}

// ── Stat bar ──────────────────────────────────────────────────────────────────

func renderStat(label string, value int, lo, hi lipgloss.Color) string {
	const barLen = 18
	bar := gradientBar(value, barLen, lo, hi)
	return fmt.Sprintf("%s %s %3d%%", labelStyle.Render(label), bar, value)
}
