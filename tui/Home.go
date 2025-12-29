package tui

import (
	"fmt"
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderProjectsSection(width int) string {
	doc := strings.Builder{}
	// --- SECTION 2: FEATURED PROJECTS (2 Cards) ---
	title := sectionTitleStyle.Render("Featured Projects")
	hintText := subtle.Align().UnsetBold().Render("Press P to view all projects")
	doc.WriteString(title + "\n")

	var projectCards []string

	limitP := 3
	if len(m.Projects) < limitP {
		limitP = len(m.Projects)
	}

	// 1. Calculate Widths
	// Card Width = (Total / 3) - Spacing
	pCardWidth := (width / 2) - 2
	if pCardWidth < 40 {
		pCardWidth = 40
	} // Safety minimum

	// Image takes fixed 32 chars. Content gets the rest.
	imageWidth := 32
	contentWidth := pCardWidth - imageWidth - 6 // -6 for padding/gap

	for i := 0; i < limitP; i++ {
		p := m.Projects[i]

		name := utils.SafeString(p, "title")
		desc := utils.SafeString(p, "description")
		githubLink := utils.SafeString(p, "githubUrl")
		liveLink := utils.SafeString(p, "liveUrl")

		// --- IMAGE HANDLING ---
		logoStr := ""
		if val, ok := p["ascii_art"].(string); ok && val != "" {
			logoStr = val
		} else if url, ok := p["imageUrl"].(string); ok && url != "" {
			logoStr = "Loading..."
		} else {
			logoStr = "   No\n  Image"
		}

		// 2. TRUNCATE IMAGE HEIGHT (The Fix)
		// Split lines, take max 10 lines, join back
		lines := strings.Split(logoStr, "\n")
		if len(lines) > 10 {
			logoStr = strings.Join(lines[:10], "\n")
		}

		// Render Image Box
		logoBox := lipgloss.NewStyle().
			Width(imageWidth).
			Align(lipgloss.Center).            // Center the ASCII horizontally in its box
			Foreground(lipgloss.Color("240")). // Grey color for image
			Render(logoStr)

		// --- CONTENT HANDLING ---
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}

		// Wrap description to fit the content column specifically
		wrappedDesc := lipgloss.NewStyle().Width(contentWidth).Render(desc)

		// Create Links Row
		links := ""
		if githubLink != "" {
			links += utils.MakeLink(" GitHub", githubLink) + "  "
		}
		if liveLink != "" {
			links += utils.MakeLink(" Live", liveLink)
		}

		contentBox := lipgloss.NewStyle().
			Width(contentWidth).
			Render(fmt.Sprintf("%s\n\n%s\n\n%s",
				highlight.Render(name),
				wrappedDesc,
				subtle.Render(links),
			))

		// 3. JOIN COLUMNS
		// Join Image + Gap + Content
		gridOf2 := lipgloss.JoinHorizontal(lipgloss.Top, logoBox, "   ", contentBox)

		// 4. RENDER CARD
		card := cardStyle.
			Width(pCardWidth).
			Height(12). // Fixed height for uniformity
			Render(gridOf2)

		projectCards = append(projectCards, card)
	}

	// Join all cards horizontally
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, projectCards...) + "\n")

	// Hint Text
	hintStyle := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Foreground(lipgloss.Color("240"))
	hintText = hintStyle.Render(hintText)

	doc.WriteString("\n" + hintText + "\n")
	return doc.String()
}

func (m Model) renderServices(width int, limitOfCards bool) string {
	doc := strings.Builder{}

	// 1. Title
	doc.WriteString(sectionTitleStyle.Render("Services Provided") + "\n\n")

	// 2. Limit to 4 Cards
	var limitS int

	if limitOfCards {
		limitS = 4
	} else {
		limitS = len(m.Services)
	}

	if len(m.Services) < limitS {
		limitS = len(m.Services)
	}

	// 3. Layout Dimensions
	// We want 2 cards per row on large screens, or 1 card on small screens
	isTwoColumn := width > 100
	var cardWidth int

	if isTwoColumn {
		cardWidth = (width / 2) - 4 // Split width minus gaps
	} else {
		cardWidth = width - 4 // Full width minus margins
	}

	var cards []string

	// 4. Style Definitions
	// The outer box for the card
	serviceCardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")). // Purple border
		Padding(1).
		MarginBottom(1).
		Width(cardWidth)

	// The Icon Box Style
	iconStyle := lipgloss.NewStyle().
		Width(6).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("205")) // Pink Icon

	// 5. Loop and Build Cards
	for i := 0; i < limitS; i++ {
		s := m.Services[i]

		// Data Extraction
		title := utils.SafeString(s, "title")
		desc := utils.SafeString(s, "description")
		price := utils.SafeString(s, "price")
		timeframe := utils.SafeString(s, "timeframe")
		category := utils.SafeString(s, "category")

		// Icon
		iconChar := utils.GetIcon(category)
		iconBox := iconStyle.Render(iconChar)

		// Content Calculation
		contentWidth := cardWidth - 10 // Card width - Icon width - Padding
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}

		// Build the Text Block
		contentBlock := lipgloss.NewStyle().Width(contentWidth).Render(fmt.Sprintf(
			"%s\n%s\n\n%s",
			highlight.Render(title),
			lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(desc),
			subtle.Render(fmt.Sprintf("%s • %s", price, timeframe)),
		))

		// Join Icon + Content
		// using Top alignment ensures icon stays at the top if text is long
		row := lipgloss.JoinHorizontal(lipgloss.Top, iconBox, "  ", contentBlock)

		// Render the full card
		cards = append(cards, serviceCardStyle.Render(row))
	}

	// 6. Layout Composition (Grid vs List)
	if isTwoColumn {
		// If we have an even number of cards, we can pair them up
		// This is a simple implementation: just JoinHorizontal pairs
		var rows []string
		for i := 0; i < len(cards); i += 2 {
			if i+1 < len(cards) {
				// Pair of cards
				rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cards[i], "  ", cards[i+1]))
			} else {
				// Single card leftover
				rows = append(rows, cards[i])
			}
		}
		doc.WriteString(lipgloss.JoinVertical(lipgloss.Left, rows...) + "\n")
	} else {
		// Standard vertical list
		doc.WriteString(lipgloss.JoinVertical(lipgloss.Left, cards...) + "\n")
	}

	// 7. Add Hint Text
	hintStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("240")).
		MarginTop(1)
	if limitOfCards {
		doc.WriteString(hintStyle.Render("For more services navigate to (S) .And for more information reach out via Contact (C) to discuss these services"))
		doc.WriteString("\n")
	} else {

		doc.WriteString(hintStyle.Render("for more information reach out via Contact (C) to discuss these services"))
	}
	doc.WriteString("\n")

	return doc.String()
}

// main rendering function for Home tab
func (m Model) renderHome(width int) string {
	doc := strings.Builder{}

	// --- SECTION 1: INTRO (Left) & CONNECT (Right) ---

	// 1. Calculate widths
	leftWidth := int(float64(width)*0.6) - 4
	rightWidth := width - leftWidth - 6

	// 2. Define Styles locally (if not global)
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginBottom(1)
	roleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	statNumber := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	statLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// --- A. LEFT COLUMN (Intro) ---

	header := fmt.Sprintf("%s \n%s",
		titleStyle.Render("Hi, I'm Tarun Nayaka R!"),
		roleStyle.Render("Freelancer   | Cloud Architect   | Full-Stack Dev  \n"),
	)

	certTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true).Underline(true)

	education := fmt.Sprintf(`
%s
%s B.Tech Computer Science And Engineering
%s PES  University
`,
		certTitle.Render("Education"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("35")).Render("🎓"),
		statLabel.Render("   2023 - 2027"),
	)

	bio := textStyle.Render(`Crafting stunning, responsive web, mobile applications 
and cloud solutions for startups and personal brands.
Specialized in Python, JS, and Cloud with 1+ years exp.`)
	br := "\n"
	stack := fmt.Sprintf(`🚀 %s
   %s Python   	%s JS/TS    	%s Go       	%s Java
   %s React    	%s Next.js  	%s Flutter  	%s Node.js
   %s Django   	%s FastAPI  	%s Express  	%s GraphQL
   %s Azure    	%s GCP      	%s AWS      	%s Docker
   %s Mongo    	%s PSQL     	%s MySQL    	%s Kafka `,
		lipgloss.NewStyle().Bold(true).Underline(true).Render("Tech Stack\n"),
		keywordStyle.Render(""), keywordStyle.Render(""), keywordStyle.Render(""), keywordStyle.Render(""),
		keywordStyle.Render(""), keywordStyle.Render(""), keywordStyle.Render(""), keywordStyle.Render(""),
		keywordStyle.Render(""), keywordStyle.Render("⚡"), keywordStyle.Render(""), keywordStyle.Render(""),
		keywordStyle.Render("ﴤ"), keywordStyle.Render(""), keywordStyle.Render(""), keywordStyle.Render(""),
		keywordStyle.Render(""), keywordStyle.Render(""), keywordStyle.Render(""), keywordStyle.Render(""),
	)
	openToWork := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true).Render("\n🟢  OPEN TO WORK")
	seeking := lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true).Render("   Seeking: Backend, Cloud\n   & Full-Stack Roles")
	statusBlock := fmt.Sprintf("%s\n%s", openToWork, seeking)

	loc := fmt.Sprintf("\n📍 %s", textStyle.Render("Bengaluru, India\n"))

	introContent := lipgloss.JoinVertical(lipgloss.Left, header, bio, education, br, stack, loc, statusBlock)

	// FIX 1: Render the Left Column into a Card!
	leftCol := cardStyle.Copy().
		Width(leftWidth).
		Render(introContent)

		// --- B. RIGHT COLUMN (Stats) ---
	statsBlock := fmt.Sprintf(`
%s %s
%s
%s %s
%s
%s %s
%s`,
		statNumber.Render("🏆 15+"), statLabel.Render("Projects Delivered"),
		statLabel.Render("   Web, Mobile,AI,ML & Cloud"),
		statNumber.Render("⏳ 1+"), statLabel.Render("Years Experience (internships)"),
		statLabel.Render("   Full-Stack & DevOps"),
		statNumber.Render("✍️ 10+"), statLabel.Render("Articles Written"),
		statLabel.Render("   Tech Blogging @ Medium"),
	)

	// 2. Certifications & Education (NEW SECTION to fill space)
	// certText := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	certsBlock := fmt.Sprintf(`
%s
%s AZ 900 Certified Azure Fundamentals
%s AZ 204 Developing Solutions for Microsoft Azure
%s Github Fundamentals
%s DP 900 Microsoft Azure Data Fundamentals
%s AI-900 Microsoft Azure AI Fundamentals
%s Google Cloud Digital Leader
`,
		certTitle.Render("Certifications"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("📜"), // Gold scroll
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("📜"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("📜"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("📜"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("📜"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("📜"),
	)

	// 3. Button-Style Links (Takes up more visual weight)
	// Define a "Button" style
	btnStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("232")). // Dark text
		Background(lipgloss.Color("250")). // Light background
		Padding(0, 1)                      // Make it chunky

	// Create styled buttons
	resumeBtn := btnStyle.Render("📄  Download Resume")
	blogBtn := btnStyle.Render("  Read My Blog   ") // Added spaces for alignment
	portfolioBtn := btnStyle.Render("  Visit Portfolio ")
	emailBtn := btnStyle.Render("  Email Me       ")

	// Wrap them in actual links
	resumeLink := utils.MakeLink(resumeBtn, "https://tarunnayaka.me/resume.pdf")
	blogLink := utils.MakeLink(blogBtn, "https://medium.com/@r.tarunnayaka25042005")
	portfolioLink := utils.MakeLink(portfolioBtn, "https://tarunnayaka.me")
	emailLink := utils.MakeLink(emailBtn, "mailto:r.tarunnayaka25042005@gmail.com")
	// Col 1: Resume & Portfolio
	col1 := lipgloss.JoinVertical(lipgloss.Left, resumeLink, "\n", portfolioLink)
	// Col 2: Blog & Email
	col2 := lipgloss.JoinVertical(lipgloss.Left, blogLink, "\n", emailLink)

	// Join Columns with a Gap
	linksBlock := lipgloss.JoinHorizontal(lipgloss.Top, col1, "   ", col2)
	// --- COMBINE RIGHT COLUMN ---
	rightContent := lipgloss.JoinVertical(lipgloss.Left,
		sectionTitleStyle.Render("At a Glance"),
		statsBlock,
		certsBlock, // Added the new Certs block here
		sectionTitleStyle.Render("Quick Links"),
		linksBlock,
	)

	rightCol := cardStyle.Copy().
		Width(rightWidth).
		Render(rightContent)

	topSection := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)
	doc.WriteString(topSection + "\n")

	//projects section ----------------------------------------------------------------------------------------
	// --- SECTION 2: FEATURED PROJECTS (3 Cards) ---
	doc.WriteString(m.renderProjectsSection(width))
	// Join all project cards horizontally

	//section for blogs
	doc.WriteString("\n")
	// --- SECTION 3: LATEST ARTICLES (3 Cards) ---
	doc.WriteString(m.renderBlogsSection(width, true))
	// Join all blog cards horizontally
	doc.WriteString("\n")

	//section 3: services offered --------------------------------------------------------------------------------
	doc.WriteString(m.renderServices(width, true))
	// --- SECTION 4: CONTACT CTA ---
	// A simple banner at the bottom
	cta := lipgloss.NewStyle().
		Width(width - 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Padding(1).
		Render("Have a project in mind? Press 'C' or Tab to visit the Contact page.")

	doc.WriteString("\n" + cta + "\n")

	return doc.String()
}
