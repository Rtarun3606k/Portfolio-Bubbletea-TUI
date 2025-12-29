package tui

import (
	"fmt"
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderBlogsSection(width int, limitOfCards bool) string {
	doc := strings.Builder{}

	// --- SECTION TITLE ---
	title := sectionTitleStyle.Render("Latest Articles")
	doc.WriteString(title + "\n\n")

	// --- 1. DETERMINE LIMIT & DATA ---
	limitB := len(m.Blogs)
	if limitOfCards == true {
		limitB = 3
	}
	if len(m.Blogs) < limitB {
		limitB = len(m.Blogs)
	}

	// --- 2. CALCULATE LAYOUT ---
	isThreeColumn := width > 120
	var cardWidth int

	if isThreeColumn {
		// Calculate width: (Total / 3) - Spacing
		cardWidth = (width / 3) - 4
	} else {
		cardWidth = width - 4
	}

	// Safety Check
	if cardWidth < 30 {
		cardWidth = 30
	}

	// Styles specific to Blog Card
	blogCardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")). // Purple/Blue border
		Padding(1).
		Width(cardWidth).
		Height(16) // Increased height slightly to fit content

	// Style for the Fallback Image Placeholder
	fallbackImageStyle := lipgloss.NewStyle().
		Width(cardWidth-4).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		Background(lipgloss.Color("236")). // Dark Grey bg
		Foreground(lipgloss.Color("245"))  // Light Grey text

	var blogCards []string

	// --- 3. LOOP & BUILD CARDS ---
	for i := 0; i < limitB; i++ {
		b := m.Blogs[i]

		// -- Data Extraction --
		id := utils.SafeString(b, "_id")
		titleVal := utils.SafeString(b, "title")
		author := utils.SafeString(b, "author")
		views := utils.SafeString(b, "views")
		dateStr := utils.SafeDate(b, "createdAt")
		id = utils.SafeID(b, "_id")

		liveLink := fmt.Sprintf("https://tarunnayaka.me/Blog/%s", id)

		// Date Formatting
		if split := strings.Split(dateStr, "T"); len(split) > 0 {
			dateStr = split[0]
		}

		// Title Truncation
		maxTitleLen := cardWidth - 6
		if len(titleVal) > maxTitleLen*2 {
			titleVal = titleVal[:(maxTitleLen*2)-3] + "..."
		}

		// --- FIX 1: SAFE IMAGE HANDLING ---
		// Instead of forcing .(string), check if it exists or use SafeString
		imgContent := utils.SafeString(b, "ascii_art")

		var imgBox string

		if imgContent != "" {
			imgBox = lipgloss.NewStyle().
				Width(cardWidth - 4).
				Align(lipgloss.Center).
				Render(imgContent)
		} else {
			// Fallback
			imgBox = fallbackImageStyle.Render("📰\nBlog Post")
		}

		// --- FIX 2: FIXED %S TYPO ---
		// Changed %S to %s
		metaText := fmt.Sprintf("%s • %s views", dateStr, views)
		metaBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(cardWidth - 4).
			Render(metaText)

		// Title Style
		titleBox := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")). // Pink title
			Width(cardWidth - 4).
			Render(titleVal)

		authorBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // Blue author
			Render("By " + author)

		link := utils.MakeLink(" Live", liveLink)
		linkHint := subtle.Render(link)

		// Combine Content
		contentBlock := lipgloss.JoinVertical(lipgloss.Left,
			metaBox,
			titleBox,
			"\n",
			authorBox,
			"\n",
			linkHint,
		)

		// Combine Image + Content + Gap
		// We add a gap between image and text using a newline or margin
		cardContent := lipgloss.JoinVertical(lipgloss.Left, imgBox, " ", contentBlock)

		// Apply Border
		finalCard := blogCardStyle.Render(cardContent)
		blogCards = append(blogCards, finalCard)
	}

	// --- 4. GRID COMPOSITION ---
	if isThreeColumn {
		var rows []string
		var currentRow []string

		// Loop through all generated cards
		for i, card := range blogCards {
			currentRow = append(currentRow, card)

			// Check if row is full (3 cards) OR if it's the last card
			if len(currentRow) == 3 || i == len(blogCards)-1 {

				// Join the 1-3 cards in this row horizontally with gaps
				renderedRow := currentRow[0]
				for k := 1; k < len(currentRow); k++ {
					renderedRow = lipgloss.JoinHorizontal(lipgloss.Top, renderedRow, "   ", currentRow[k])
				}

				// Center this row within the full width and add to our list of rows
				rows = append(rows, lipgloss.PlaceHorizontal(width, lipgloss.Center, renderedRow))

				// Clear current row for the next batch
				currentRow = []string{}
			}
		}

		// Now join all the rows Vertically
		// Add a vertical gap ("\n") between rows if you like
		finalGrid := lipgloss.JoinVertical(lipgloss.Left, rows...)
		doc.WriteString(finalGrid)

	} else {
		// Vertical Stack for small screens (Mobile)
		doc.WriteString(lipgloss.JoinVertical(lipgloss.Center, blogCards...))
	}

	doc.WriteString("\n")

	// --- 5. HINT TEXT ---
	hintStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	if limitOfCards {
		doc.WriteString(hintStyle.Render("Press (B) to view all articles • Enter to read"))
	} else {
		doc.WriteString(hintStyle.Render("Press (Esc) to return home • Enter to read"))
	}

	doc.WriteString("\n")

	return doc.String()
}
