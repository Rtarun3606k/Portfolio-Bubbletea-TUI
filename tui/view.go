package tui

import (
	"fmt"
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	subtle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// Navigation Bar Styles
	activeTabBorder = lipgloss.Border{
		Top: "─", Bottom: " ", Left: "│", Right: "│",
		TopLeft: "╭", TopRight: "╮", BottomLeft: "┘", BottomRight: "└",
	}
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), true)
	activeTabStyle = tabStyle.Copy().Border(activeTabBorder, true).Foreground(lipgloss.Color("63"))
)

func (m Model) generateConetnt(width int) string {

	// 1. Calculate Content Dimensions (e.g., 80% of screen width)
	contentWidth := int(float64(m.Width) * 0.8)
	if contentWidth < 40 {
		contentWidth = 40 // Safe minimum width
	}
	doc := strings.Builder{}

	switch m.ActiveTab {
	case 0:
		return m.renderHome(width)

	case 1: // Projects
		return m.renderProject(width)

	case 2:
		return m.renderPosition(contentWidth)

	case 3: // Contact
		doc.WriteString("Send me a message:\n\n")
		doc.WriteString(m.EmailInput.View() + "\n\n")
		doc.WriteString(m.MsgInput.View())
	}

	return doc.String()

}

func (m Model) View() string {
	// 1. Loading Screen
	if m.Loading {
		return lipgloss.Place(
			m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			fmt.Sprintf("%s Loading Data...", m.Spinner.View()),
		)
	}

	// 2. BUILD HEADER (Logo + Gap + Tabs)
	// Logo Style
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("235")).
		Bold(true).
		Padding(0, 1).
		MarginRight(1).
		MarginLeft(3).
		MarginTop(2).
		SetString("TARUN NAYAKA R")

	logo := logoStyle.Render()

	// Tabs Style
	tabs := []string{"  Home (H)", "  Projects (P)", "  Experience (E)", "  Contact (C)"}
	var renderedTabs []string

	for i, t := range tabs {
		if m.ActiveTab == i {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(t))
		}
	}
	tabsBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Calculate Gap
	gapWidth := m.Width - lipgloss.Width(logo) - lipgloss.Width(tabsBlock) - 4 // -4 for extra safety margin
	if gapWidth < 0 {
		gapWidth = 0
	}
	gap := strings.Repeat(" ", gapWidth)

	// Combine into Header
	header := lipgloss.JoinHorizontal(lipgloss.Top, logo, gap, tabsBlock)
	// Add some padding below the header
	header = lipgloss.NewStyle().MarginBottom(1).Render(header)

	// 3. BUILD VIEWPORT (Content)
	viewportContent := lipgloss.NewStyle().
		Width(m.Width).
		Align(lipgloss.Center).
		Render(m.Viewport.View())

		// 4. BUILD FOOTER (Help Text)
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("use ← → or Tab to navigate • j/k to scroll • q to quit")

	//  Social Links
	ghIcon := highlight.Render("  GitHub")
	liIcon := highlight.Render("  LinkedIn")
	webIcon := highlight.Render("  Portfolio") // or use  for Desktop

	github := utils.MakeLink(ghIcon, "https://github.com/tarunNayaka")
	linkedin := utils.MakeLink(liIcon, "https://linkedin.com/in/tarun")
	portFolio := utils.MakeLink(webIcon, "https://tarunnayaka.me")

	// Join them with some spacing
	socials := lipgloss.JoinHorizontal(lipgloss.Top, github, "   ", linkedin, "   ", portFolio)

	//  Combine Help + Socials into one block
	footerContent := lipgloss.JoinVertical(lipgloss.Center, helpText, " \n", socials)

	//  Render the full footer container
	footer := lipgloss.NewStyle().
		Width(m.Width).
		Align(lipgloss.Center).
		PaddingTop(1). // Space between content and footer
		Render(footerContent)

	//  STACK EVERYTHING VERTICALLY
	return lipgloss.JoinVertical(lipgloss.Top,
		header,
		viewportContent,
		footer,
	)
}
