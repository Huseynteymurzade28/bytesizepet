package main

import "github.com/charmbracelet/lipgloss"

// Pastel cozy palette — Nord / cottagecore inspired.
const (
	colorLavender = lipgloss.Color("#c4b5e0")
	colorRose     = lipgloss.Color("#f0a0c0")
	colorSage     = lipgloss.Color("#a0c4a0")
	colorSalmon   = lipgloss.Color("#e8907a")
	colorSunbeam  = lipgloss.Color("#f5d070")
	colorSkyBlue  = lipgloss.Color("#78c8e8")
	colorMint     = lipgloss.Color("#78d8a8")
	colorCream    = lipgloss.Color("#f0e8d8")
	colorMuted    = lipgloss.Color("#8888aa")
	colorDimBar   = lipgloss.Color("#3a3858")
	colorPlum     = lipgloss.Color("#7a5c98")
	colorAmber    = lipgloss.Color("#f0b060")
	colorCoral    = lipgloss.Color("#e07868")
	colorStarDust = lipgloss.Color("#b8a8d8")
	colorHabitat  = lipgloss.Color("#7060a0")
)

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorLavender).
			Padding(1, 3)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRose).
			Align(lipgloss.Center).
			Width(44)

	labelStyle = lipgloss.NewStyle().
			Width(11).
			Foreground(colorMuted)

	barEmptyStyle = lipgloss.NewStyle().
			Foreground(colorDimBar)

	statusStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(colorCream)

	statusBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPlum).
			Foreground(colorCream).
			Italic(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorStarDust)

	dividerStyle = lipgloss.NewStyle().
			Foreground(colorPlum)

	habitatStyle = lipgloss.NewStyle().
			Foreground(colorHabitat)

	petAreaStyle = lipgloss.NewStyle().
			Width(44).
			Align(lipgloss.Center)

	// Character select cards
	charBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(1, 1).
			Width(14).
			Align(lipgloss.Center)

	charBoxSelectedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorRose).
				Padding(1, 1).
				Width(14).
				Align(lipgloss.Center)

	// Minigame
	mgTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSunbeam)

	mgRPSTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSkyBlue)

	mgMathTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorMint)

	reactionWaitStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Italic(true)

	reactionReadyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorSage)

	resultStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRose)

	mathQuestionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorSunbeam).
				Align(lipgloss.Center)
)
