package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Root view ─────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.quitting {
		farewell := lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Render("\n  Hoşça kal! " + m.pet.name + " seni özleyecek~ 🐾\n\n")
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
	title := titleStyle.Render("✨  Karakterini Seç  ✨")
	sub := helpStyle.Render("  ← / → ile seç,  Enter ile onayla")

	boxes := make([]string, len(AllChars))
	for i, ch := range AllChars {
		art := lipgloss.NewStyle().Foreground(ch.Colors.Normal).Render(ch.Art.Normal)
		label := lipgloss.NewStyle().Bold(i == m.selectedChar).Render(ch.Name + "\n" + ch.Emoji + " " + ch.Species)
		inner := art + "\n\n" + label

		if i == m.selectedChar {
			boxes[i] = charBoxSelectedStyle.Render(inner)
		} else {
			boxes[i] = charBoxStyle.Render(inner)
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, boxes...)

	content := strings.Join([]string{title, sub, "", row}, "\n")
	return borderStyle.Render(content)
}

// ── Main screen ───────────────────────────────────────────────────────────────

func viewMain(m model) string {
	p := m.pet

	title := titleStyle.Render(fmt.Sprintf("✨  %s  ✨", p.name))
	art := renderPetArt(p)

	indicator := renderStatusIndicator(p)
	infoRow := fmt.Sprintf("Durum: %s    Yaş: %d tık", indicator, p.age)

	stats := strings.Join([]string{
		renderStat("Açlık  :", p.hunger, hungerColor(p.hunger)),
		renderStat("Mutluluk:", p.happiness, happinessColor(p.happiness)),
		renderStat("Enerji  :", p.energy, energyColor(p.energy)),
	}, "\n")

	status := statusStyle.Render("» " + p.statusMsg)
	divider := dividerStyle.Render(strings.Repeat("─", 40))
	help := helpStyle.Render("[f] Besle  [p] Oyna  [s] Uyu  [1] Tahmin  [2] Tepki  [q] Çık")

	content := strings.Join([]string{
		title, "",
		art, "",
		infoRow, "",
		stats, "",
		status,
		divider,
		help,
	}, "\n")

	return borderStyle.Render(content)
}

// ── Minigame screen ───────────────────────────────────────────────────────────

func viewMinigame(m model) string {
	title := titleStyle.Render(fmt.Sprintf("✨  %s  ✨", m.pet.name))
	art := renderPetArt(m.pet)

	var gameSection string
	switch {
	case m.mg.phase == mgGuessPlaying || m.mg.phase == mgGuessDone:
		gameSection = viewGuessGame(m)
	default:
		gameSection = viewReactionGame(m)
	}

	help := helpStyle.Render("[ESC] Ana ekrana dön")

	content := strings.Join([]string{
		title, "",
		art, "",
		gameSection, "",
		help,
	}, "\n")

	return borderStyle.Render(content)
}

func viewGuessGame(m model) string {
	header := mgTitleStyle.Render("🎮  TAHMİN OYUNU")

	if m.mg.phase == mgGuessDone {
		return strings.Join([]string{
			header, "",
			resultStyle.Render(m.mg.resultMsg), "",
			helpStyle.Render("[ENTER] Devam"),
		}, "\n")
	}

	hint := statusStyle.Render(m.mg.hint)
	keys := helpStyle.Render("[1-9] Tahmin et")

	return strings.Join([]string{header, "", hint, keys}, "\n")
}

func viewReactionGame(m model) string {
	header := mgTitleStyle.Render("⚡  TEPKİ OYUNU")

	var body string
	switch m.mg.phase {
	case mgReactionWait:
		body = reactionWaitStyle.Render("Hazırlan... sinyal bekleniyor...") +
			"\n\n" + helpStyle.Render("[SPACE] beklet — erken basma!")
	case mgReactionReady:
		body = reactionReadyStyle.Render("  ★  ŞİMDİ!  ★  ") +
			"\n\n" + helpStyle.Render("[SPACE] HEMEN BAS!")
	case mgReactionDone:
		body = resultStyle.Render(m.mg.resultMsg) +
			"\n\n" + helpStyle.Render("[ENTER] Devam")
	}

	return strings.Join([]string{header, "", body}, "\n")
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

	// Append zzz visually when sleeping without altering the art string itself.
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
		return col("147", "💤 Uyuyor")
	case p.hunger >= 75:
		return col("196", "🍖 Aç!")
	case p.happiness <= 20:
		return col("196", "😢 Mutsuz!")
	case p.energy <= 20:
		return col("208", "😪 Yorgun!")
	default:
		return col("82", "😊 Mutlu")
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
