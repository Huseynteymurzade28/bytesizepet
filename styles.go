package main

import "github.com/charmbracelet/lipgloss"

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("219")).
			Padding(1, 4)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213")).
			Align(lipgloss.Center).
			Width(44)

	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("147")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Width(11).
			Foreground(lipgloss.Color("183"))

	barEmptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	statusStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("222"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("219"))

	charBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(1, 1).
			Width(14).
			Align(lipgloss.Center)

	charBoxSelectedStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("213")).
				Padding(1, 1).
				Width(14).
				Align(lipgloss.Center)

	mgTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226"))

	mgRPSTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("159"))

	mgMathTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("183"))

	reactionWaitStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Italic(true)

	reactionReadyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("82"))

	resultStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213"))

	mathQuestionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("226")).
				Align(lipgloss.Center)
)
