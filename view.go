package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Root view ─────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.quitting {
		msg := "Come back soon~ 🐾"
		if m.pet.name != "" {
			msg = m.pet.name + " will miss you~ 🐾"
		}
		farewell := lipgloss.NewStyle().Foreground(colorRose).
			Render("\n  Goodbye! " + msg + "\n\n")
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, farewell)
	}

	var content string
	switch m.screen {
	case screenCharSelect:
		content = viewCharSelect(m)
	case screenMain:
		content = viewMain(m)
	case screenMinigame:
		content = viewMinigame(m)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// ── Character select ──────────────────────────────────────────────────────────

func viewCharSelect(m model) string {
	title := titleStyle.Render("✨  Choose Your Pet  ✨")
	sub := helpStyle.Render("  ← / → / ↑ / ↓  to select,  Enter to confirm")

	boxes := make([]string, len(AllChars))
	for i, ch := range AllChars {
		art := lipgloss.NewStyle().Foreground(ch.Colors.Normal).Render(ch.Art.Normal)
		nameStr := ch.Name
		if i == m.selectedChar {
			nameStr = lipgloss.NewStyle().Bold(true).Foreground(colorRose).Render("✦ " + nameStr + " ✦")
		}
		label := nameStr + "\n" + ch.Emoji + " " + ch.Species
		inner := art + "\n\n" + label

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
	return borderStyle.Render(content)
}

// ── Main screen ───────────────────────────────────────────────────────────────

func viewMain(m model) string {
	p := m.pet

	title := titleStyle.Render(fmt.Sprintf("✨  %s  %s  ✨", p.name, AllChars[p.charIdx].Emoji))

	centeredPet := petAreaStyle.Render(renderAnimatedPet(m))

	indicator := renderStatusIndicator(p)
	infoRow := fmt.Sprintf("Status: %s    Age: %d ticks", indicator, p.age)

	stats := strings.Join([]string{
		renderStat("Hunger   :", p.hunger, hungerColor(p.hunger)),
		renderStat("Happiness:", p.happiness, happinessColor(p.happiness)),
		renderStat("Energy   :", p.energy, energyColor(p.energy)),
	}, "\n")

	statusBox := statusBoxStyle.Render(statusIcon(p) + " " + p.statusMsg)

	divider := dividerStyle.Render(strings.Repeat("─", 40))
	controls1 := renderControls(
		[]string{"f", "p", "s", "q"},
		[]string{"Feed", "Play", "Sleep", "Quit"},
	)
	controls2 := renderControls(
		[]string{"1", "2", "3", "4"},
		[]string{"Guess", "React", "RPS", "Math"},
	)

	content := strings.Join([]string{
		title, "",
		centeredPet, "",
		infoRow, "",
		stats, "",
		statusBox,
		divider,
		controls1,
		controls2,
	}, "\n")

	return borderStyle.Render(content)
}

// ── Minigame screen ───────────────────────────────────────────────────────────

func viewMinigame(m model) string {
	title := titleStyle.Render(fmt.Sprintf("✨  %s  %s  ✨", m.pet.name, AllChars[m.pet.charIdx].Emoji))
	centeredPet := petAreaStyle.Render(renderAnimatedPet(m))

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
		title, "",
		centeredPet, "",
		gameSection, "",
		help,
	}, "\n")

	return borderStyle.Render(content)
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

// ── Animated pet renderer ─────────────────────────────────────────────────────

func renderAnimatedPet(m model) string {
	p := m.pet
	ch := AllChars[p.charIdx]

	var art string
	var color lipgloss.Color

	isDead := p.hunger >= 90 || (p.happiness <= 5 && p.energy <= 0)

	switch {
	case isDead:
		art, color = ch.Art.Dead, "#707090"
	case p.sleeping:
		art, color = ch.Art.Sleep, ch.Colors.Sleep
	case p.energy <= 20:
		art, color = ch.Art.Tired, ch.Colors.Tired
	case p.hunger >= 75:
		if m.animFrame == 1 {
			art = shakeArt(ch.Art.Hungry)
		} else {
			art = ch.Art.Hungry
		}
		color = ch.Colors.Hungry
	case p.happiness >= 80:
		art, color = ch.Art.Happy, ch.Colors.Happy
	default:
		if m.animFrame == 0 {
			art = ch.Art.Normal
		} else {
			art = ch.Art.NormalB
		}
		color = ch.Colors.Normal
	}

	rendered := lipgloss.NewStyle().Foreground(color).Render(art)

	// Animated sleep z's float above the head
	if p.sleeping {
		zFrames := [4]string{"   ", "z  ", "zZ ", "zZz"}
		zStr := lipgloss.NewStyle().Foreground(lipgloss.Color("#b4a8e8")).Italic(true).
			Render(" " + zFrames[m.sleepFrame])
		lines := strings.Split(rendered, "\n")
		lines[0] += zStr
		rendered = strings.Join(lines, "\n")
	}

	// Critical-hunger sweat drop on alternate frames
	if p.hunger >= 85 && !p.sleeping && !isDead && m.animFrame == 1 {
		drop := lipgloss.NewStyle().Foreground(lipgloss.Color("#78c8e8")).Render("~")
		lines := strings.Split(rendered, "\n")
		lines[0] = drop + lines[0]
		rendered = strings.Join(lines, "\n")
	}

	// Habitat — cycles between frame A and B
	hab := ch.HabitatA
	if m.envFrame == 1 {
		hab = ch.HabitatB
	}
	rendered += "\n" + habitatStyle.Render(hab)

	// Action animation line (always present to prevent layout shift)
	rendered += "\n" + renderActionAnim(m.actionAnim, m.actionFrame)

	return rendered
}

// shakeArt offsets alternating lines by one space to simulate trembling.
func shakeArt(art string) string {
	lines := strings.Split(art, "\n")
	for i := range lines {
		if i%2 == 0 {
			lines[i] = " " + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}

// renderActionAnim returns a one-line animation for feed/play actions.
// Returns spaces when inactive so the layout height stays constant.
func renderActionAnim(anim, frame int) string {
	if frame >= 4 {
		frame = 3
	}
	feedFrames := [4]string{
		">°>            ",
		"  >°>          ",
		"    >°>        ",
		"      nom! ♪   ",
	}
	playFrames := [4]string{
		"  ♡             ",
		" ♡  ♡           ",
		"♡  ♡  ♡         ",
		"  ·  ·  ·       ",
	}
	switch anim {
	case 1:
		return lipgloss.NewStyle().Foreground(colorSalmon).Render(feedFrames[frame])
	case 2:
		return lipgloss.NewStyle().Foreground(colorRose).Render(playFrames[frame])
	}
	return "                "
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

// ── Status indicator ──────────────────────────────────────────────────────────

func renderStatusIndicator(p petModel) string {
	col := func(c lipgloss.Color, t string) string {
		return lipgloss.NewStyle().Foreground(c).Render(t)
	}
	switch {
	case p.sleeping:
		return col(colorStarDust, "💤 Sleeping")
	case p.hunger >= 75:
		return col(colorCoral, "🍖 Hungry!")
	case p.happiness <= 20:
		return col(colorCoral, "😢 Sad!")
	case p.energy <= 20:
		return col(colorAmber, "😪 Tired!")
	default:
		return col(colorSage, "😊 Happy")
	}
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

func renderStat(label string, value int, fillColor lipgloss.Color) string {
	const barLen = 18
	filled := value * barLen / 100
	bar := lipgloss.NewStyle().Foreground(fillColor).Render(strings.Repeat("█", filled)) +
		barEmptyStyle.Render(strings.Repeat("░", barLen-filled))
	return fmt.Sprintf("%s %s %3d%%", labelStyle.Render(label), bar, value)
}

// ── Color helpers ─────────────────────────────────────────────────────────────

func hungerColor(v int) lipgloss.Color {
	switch {
	case v >= 70:
		return colorCoral
	case v >= 40:
		return colorAmber
	default:
		return colorSalmon
	}
}

func happinessColor(v int) lipgloss.Color {
	switch {
	case v <= 30:
		return colorCoral
	case v <= 60:
		return colorAmber
	default:
		return colorSunbeam
	}
}

func energyColor(v int) lipgloss.Color {
	switch {
	case v <= 20:
		return colorCoral
	case v <= 50:
		return colorAmber
	default:
		return colorSkyBlue
	}
}
