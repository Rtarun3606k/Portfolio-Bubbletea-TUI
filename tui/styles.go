package tui

import "github.com/charmbracelet/lipgloss"

var (

	// Card Styles
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			MarginRight(1) // Gap between cards

	// Section Titles
	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				Underline(true).
				MarginBottom(1).
				MarginTop(2) // Space before new section
)
