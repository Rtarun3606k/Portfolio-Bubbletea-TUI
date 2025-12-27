package tui

import (
	"fmt"
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m Model) renderPosition(width int) string {
	doc := strings.Builder{}

	// --- 1. Define Layout Dimensions ---
	// We split the card into: Left (Logo) + Right (Content)
	// Total available width inside the border:
	cardWidth := width - 4                    // Account for border/padding
	logoWidth := 32                           // Fixed width for ASCII art column
	contentWidth := cardWidth - logoWidth - 3 // Remaining space (-3 for gap)

	if contentWidth < 20 {
		contentWidth = 20 // Safety minimum
	}

	// --- 2. Define Local Styles ---
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")). // Purple Border
		Padding(1, 1).
		MarginBottom(1).
		Width(cardWidth)

	// Style for the Meta row (Internship • IND • Remote)
	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")). // Dark Grey
		Italic(true)

	// --- 3. Iterate Experience ---
	for _, e := range m.Experience {

		// A. PREPARE DATA
		role := utils.SafeString(e, "jobTitle")
		company := utils.SafeString(e, "companyName")
		empType := utils.SafeString(e, "employmentType")
		location := utils.SafeString(e, "location")

		// Dates
		startDate := utils.SafeString(e, "startDate")
		if len(startDate) > 10 {
			startDate = startDate[:10]
		}

		endDate := utils.SafeString(e, "endDate")
		if len(endDate) > 10 {
			endDate = endDate[:10]
		}

		// Booleans
		isCurrent := false
		if val, ok := e["isCurrent"].(bool); ok {
			isCurrent = val
		}
		if isCurrent {
			endDate = "Present"
		}

		isRemote := false
		if val, ok := e["isRemote"].(bool); ok {
			isRemote = val
		}

		remoteStr := "On-site"
		if isRemote {
			remoteStr = "Remote"
		}

		// B. BUILD LEFT COLUMN (Logo)
		// We use the ASCII art we generated earlier
		logoStr := "   No\n  Image"
		if val, ok := e["ascii_art"].(string); ok && val != "" {
			logoStr = val
		}

		// Render Left Column
		leftCol := lipgloss.NewStyle().
			Width(logoWidth).
			Align(lipgloss.Center).
			Render(logoStr)

		// C. BUILD RIGHT COLUMN (Info)

		// 1. Header: Role @ Company
		header := fmt.Sprintf("%s @ %s", highlight.Render(role), highlight.Render(company))

		// 2. Meta: Internship • IND • Remote
		metaInfo := fmt.Sprintf("%s • %s • %s", empType, location, remoteStr)

		// 3. Date Row
		dateRow := fmt.Sprintf("🗓  %s - %s", startDate, endDate)

		// 4. Responsibilities (Wrapped to Content Width!)
		var resBuilder strings.Builder

		// Helper to process list
		processList := func(items []interface{}) {
			for _, r := range items {
				if str, ok := r.(string); ok {
					// CRITICAL: Wrap text to 'contentWidth', not full screen width
					wrapped := lipgloss.NewStyle().Width(contentWidth - 2).Render(str)
					resBuilder.WriteString(fmt.Sprintf("• %s\n", wrapped))
				}
			}
		}

		if rawResp, ok := e["responsibilities"].(bson.A); ok {
			processList(rawResp)
		} else if rawResp, ok := e["responsibilities"].([]interface{}); ok {
			processList(rawResp)
		}

		// Assemble Right Column
		rightBlock := lipgloss.JoinVertical(lipgloss.Left,
			header,
			metaStyle.Render(metaInfo),
			subtle.Render(dateRow),
			"\n", // Spacer
			resBuilder.String(),
		)

		// D. JOIN COLUMNS
		// Put Logo (Left) and Content (Right) side-by-side
		row := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "   ", rightBlock)

		// E. RENDER CARD
		doc.WriteString(cardStyle.Render(row) + "\n")
	}

	return doc.String()
}
