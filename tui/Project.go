package tui

import (
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m Model) renderProject(width int) string {
	doc := strings.Builder{}

	// --- 1. Layout Dimensions ---
	// Split: Left (Image) + Right (Content)
	cardWidth := width - 4
	imageWidth := 32 // Matches the width used in generation (30) + padding
	contentWidth := cardWidth - imageWidth - 3

	if contentWidth < 20 {
		contentWidth = 20
	}

	// --- 2. Styles ---
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")). // Pink border for Projects
		Padding(1, 1).
		MarginBottom(1).
		Width(cardWidth)

	// Badge Style
	featuredBadge := lipgloss.NewStyle().
		Foreground(lipgloss.Color("228")). // Yellow/Gold
		Background(lipgloss.Color("63")).  // Purple bg
		Bold(true).
		Padding(0, 1).
		SetString("★ FEATURED")

	// Tag Style (e.g., [Python])
	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("123")). // Cyan text
		Background(lipgloss.Color("237")). // Dark grey bg
		Padding(0, 1).
		MarginRight(1)

	// --- 3. Iterate Projects ---
	for _, p := range m.Projects {

		// A. EXTRACT DATA
		title := utils.SafeString(p, "title")
		desc := utils.SafeString(p, "description")
		github := utils.SafeString(p, "githubUrl")
		live := utils.SafeString(p, "liveUrl")

		// Featured Status
		isFeatured := false
		if val, ok := p["featured"].(bool); ok {
			isFeatured = val
		}
		//https://walkez.blob.core.windows.net/projects/1757574625363-Screenshotfrom2025-08-0114-25-42.png

		// Tags (Handle bson.A)
		var tagList []string
		if rawTags, ok := p["tags"].(bson.A); ok {
			for _, t := range rawTags {
				if str, ok := t.(string); ok {
					tagList = append(tagList, str)
				}
			}
		}

		// B. BUILD LEFT COLUMN (ASCII Art)
	fallbackImageStyle := lipgloss.NewStyle().
		Width(cardWidth-4).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		Background(lipgloss.Color("236")). // Dark Grey bg
		Foreground(lipgloss.Color("245"))  // Light Grey text


		imgContent := utils.SafeString(p, "ascii_art")

		var imgBox string

		if imgContent != "" {
			imgBox = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Render(imgContent)
		} else {
			// Fallback
			imgBox = fallbackImageStyle.Render("📰\nBlog Post")
		}

		// Vertically center the image roughly if description is long
		leftCol := lipgloss.NewStyle().
			Width(imageWidth).
			Align(lipgloss.Center).
			Render(imgBox)

		// C. BUILD RIGHT COLUMN

		// 1. Title Row (Title + Featured Badge)
		titleRow := highlight.Render(title)
		if isFeatured {
			titleRow += "  " + featuredBadge.String()
		}

		// 2. Links Row
		var links []string
		if github != "" {
			links = append(links, utils.MakeLink("  Source Code", github))
		}
		if live != "" {
			links = append(links, utils.MakeLink("🔗  Live Demo", live))
		}
		linkRow := strings.Join(links, "   ") // Spacing between links

		// 3. Tags Row
		var styledTags []string
		for _, t := range tagList {
			styledTags = append(styledTags, tagStyle.Render(t))
		}
		tagsRow := lipgloss.JoinHorizontal(lipgloss.Top, styledTags...)

		// 4. Description (Wrapped)
		wrappedDesc := lipgloss.NewStyle().
			Width(contentWidth).
			Foreground(lipgloss.Color("252")). // White/Grey text
			Render(desc)

		// Assemble Right Stack
		// Order: Title -> Links -> Tags -> Description
		rightStack := lipgloss.JoinVertical(lipgloss.Left,
			titleRow,
			linkRow,
			"\n", // Spacer
			tagsRow,
			"\n", // Spacer
			wrappedDesc,
		)

		// D. COMBINE COLUMNS
		// Image |  Details
		row := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "   ", rightStack)

		// E. RENDER CARD
		doc.WriteString(cardStyle.Render(row) + "\n")
	}

	return doc.String()
}
