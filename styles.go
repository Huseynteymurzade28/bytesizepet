package main

import "github.com/charmbracelet/lipgloss"

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 4)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213")).
			Align(lipgloss.Center).
			Width(40)

	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Width(11).
			Foreground(lipgloss.Color("245"))

	barEmptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	statusStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("222"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("237"))

	charBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(1, 2).
			Width(14).
			Align(lipgloss.Center)

	charBoxSelectedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("213")).
				Padding(1, 2).
				Width(14).
				Align(lipgloss.Center)

	mgTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226"))

	reactionWaitStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Italic(true)

	reactionReadyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("82"))

	resultStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213"))
)
