package tui

import (
	"fmt"
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	case 0: // Home
		// doc.WriteString(highlight.Render("WELCOME TO MY TERMINAL") + "\n\n")
		//
		// for _, s := range m.Services {
		// 	title := utils.SafeString(s, "title")
		// 	price := utils.SafeString(s, "price")
		//
		// 	// Wrap title to fit container
		// 	wrappedTitle := lipgloss.NewStyle().Width(contentWidth).Render(title)
		// 	doc.WriteString(fmt.Sprintf("• %s  %s\n", wrappedTitle, subtle.Render(price)))
		// }
		// githubicon := highlight.Render("  GitHub")
		// Linkedinicon := highlight.Render("  LinkedIn")
		//
		// // 2. Wrap that styled string in the link
		// github := utils.MakeLink(githubicon, "https://github.com/tarunNayaka")
		// linkedin := utils.MakeLink(Linkedinicon, "https://github.com/tarunNayaka")
		//
		// doc.WriteString("\n" + github + "   ")
		// doc.WriteString(" " + linkedin + "\n")
		return m.renderHome(width)

	case 1: // Projects
		for _, p := range m.Projects {
			name := utils.SafeString(p, "title")
			desc := utils.SafeString(p, "description")
			github := utils.SafeString(p, "githubUrl")
			live := utils.SafeString(p, "liveUrl")

			if github != "" {
				name += " " + utils.MakeLink("", github)
			}
			if live != "" {
				name += " " + utils.MakeLink("🔗", live)
			}

			// Wrap description to 100% of the Container Width
			wrapperDesc := lipgloss.NewStyle().Width(contentWidth).Render(desc)

			doc.WriteString(fmt.Sprintf("%s\n%s\n%s\n%s\n\n", highlight.Render(name), wrapperDesc, subtle.Render(github), subtle.Render(live)))
		}

	case 2: // Experience
		for _, e := range m.Experience {
			// 1. Build the Responsibilities String safely
			var res string

			// Helper to wrap list items
			wrapList := func(items []interface{}) {
				for _, r := range items {
					if str, ok := r.(string); ok {
						// Wrap text to fit container width (minus 2 chars for bullet)
						wrapped := lipgloss.NewStyle().Width(contentWidth - 2).Render(str)
						res += fmt.Sprintf("• %s\n", wrapped)
					}
				}
			}

			if rawResp, ok := e["responsibilities"].(bson.A); ok {
				// Convert bson.A to []interface{} logic
				// (Since bson.A is literally []interface{}, we can loop directly)
				for _, r := range rawResp {
					if str, ok := r.(string); ok {
						wrapped := lipgloss.NewStyle().Width(contentWidth - 2).Render(str)
						res += fmt.Sprintf("• %s\n", wrapped)
					}
				}
			} else if rawResp, ok := e["responsibilities"].([]interface{}); ok {
				wrapList(rawResp)
			}

			// 2. Safe Field Access
			role := utils.SafeString(e, "jobTitle")
			company := utils.SafeString(e, "companyName")
			startDate := utils.SafeString(e, "startDate")
			if len(startDate) > 10 {
				startDate = startDate[:10]
			}
			endDate := utils.SafeString(e, "endDate")
			if len(endDate) > 10 {
				endDate = endDate[:10]
			}

			isCurrent := false
			if val, ok := e["isCurrent"].(bool); ok {
				isCurrent = val
			}
			if isCurrent {
				endDate = "Present"
			}

			// 3. Render Output
			doc.WriteString(highlight.Render(fmt.Sprintf("%s @ %s", role, company)) + "\n")
			doc.WriteString(subtle.Render(fmt.Sprintf("%s - %s", startDate, endDate)) + "\n\n")

			if res != "" {
				doc.WriteString(res + "\n")
			}
		}

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
