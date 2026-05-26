package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Root view ─────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.quitting {
		var msg string
		if m.pet.name != "" {
			msg = "Goodbye! " + m.pet.name + " will miss you~ 🐾"
		} else {
			msg = "Goodbye! Come back soon~ 🐾"
		}
		farewell := lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Render("\n  " + msg + "\n\n")
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
			nameStr = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("213")).Render(nameStr)
		}
		label := nameStr + "\n" + ch.Emoji + " " + ch.Species
		inner := art + "\n\n" + label

		if i == m.selectedChar {
			boxes[i] = charBoxSelectedStyle.Render(inner)
		} else {
			boxes[i] = charBoxStyle.Render(inner)
		}
	}

	// Layout in rows of charSelectCols
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
	art := renderPetArt(p)

	indicator := renderStatusIndicator(p)
	infoRow := fmt.Sprintf("Status: %s    Age: %d ticks", indicator, p.age)

	stats := strings.Join([]string{
		renderStat("Hunger   :", p.hunger, hungerColor(p.hunger)),
		renderStat("Happiness:", p.happiness, happinessColor(p.happiness)),
		renderStat("Energy   :", p.energy, energyColor(p.energy)),
	}, "\n")

	status := statusStyle.Render("» " + p.statusMsg)
	divider := dividerStyle.Render(strings.Repeat("─", 40))
	help1 := helpStyle.Render("[f] Feed  [p] Play  [s] Sleep  [q] Quit")
	help2 := helpStyle.Render("[1] Guess  [2] React  [3] RPS  [4] Math")

	content := strings.Join([]string{
		title, "",
		art, "",
		infoRow, "",
		stats, "",
		status,
		divider,
		help1,
		help2,
	}, "\n")

	return borderStyle.Render(content)
}

// ── Minigame screen ───────────────────────────────────────────────────────────

func viewMinigame(m model) string {
	title := titleStyle.Render(fmt.Sprintf("✨  %s  %s  ✨", m.pet.name, AllChars[m.pet.charIdx].Emoji))
	art := renderPetArt(m.pet)

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
		art, "",
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
		m.mg.mathOptions[0],
		m.mg.mathOptions[1],
		m.mg.mathOptions[2],
	))

	return strings.Join([]string{header, "", question, "", opts}, "\n")
}

// ── Pet art ───────────────────────────────────────────────────────────────────

func renderPetArt(p petModel) string {
	ch := AllChars[p.charIdx]

	var art string
	var color lipgloss.Color

	switch {
	case p.hunger >= 90 || (p.happiness <= 5 && p.energy <= 0):
		art, color = ch.Art.Dead, "240"
	case p.sleeping:
		art, color = ch.Art.Sleep, ch.Colors.Sleep
	case p.energy <= 20:
		art, color = ch.Art.Tired, ch.Colors.Tired
	case p.hunger >= 75:
		art, color = ch.Art.Hungry, ch.Colors.Hungry
	case p.happiness >= 80:
		art, color = ch.Art.Happy, ch.Colors.Happy
	default:
		art, color = ch.Art.Normal, ch.Colors.Normal
	}

	rendered := lipgloss.NewStyle().Foreground(color).Render(art)

	if p.sleeping {
		zzz := lipgloss.NewStyle().Foreground(lipgloss.Color("147")).Italic(true).Render(" zZz")
		lines := strings.Split(rendered, "\n")
		lines[0] += zzz
		return strings.Join(lines, "\n")
	}

	return rendered
}

// ── Status indicator ──────────────────────────────────────────────────────────

func renderStatusIndicator(p petModel) string {
	col := func(c, t string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render(t)
	}
	switch {
	case p.sleeping:
		return col("147", "💤 Sleeping")
	case p.hunger >= 75:
		return col("196", "🍖 Hungry!")
	case p.happiness <= 20:
		return col("196", "😢 Sad!")
	case p.energy <= 20:
		return col("208", "😪 Tired!")
	default:
		return col("82", "😊 Happy")
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
		return "196"
	case v >= 40:
		return "214"
	default:
		return "82"
	}
}

func happinessColor(v int) lipgloss.Color {
	switch {
	case v <= 30:
		return "196"
	case v <= 60:
		return "214"
	default:
		return "82"
	}
}

func energyColor(v int) lipgloss.Color {
	switch {
	case v <= 20:
		return "196"
	case v <= 50:
		return "214"
	default:
		return "75"
	}
}
