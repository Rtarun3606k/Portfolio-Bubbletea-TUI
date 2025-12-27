package tui

import (
	"fmt"
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
	doc.WriteString(sectionTitleStyle.Render("Featured Projects") + "\n")

	var projectCards []string
	// Limit to first 3 projects (safety check)
	limitP := 3
	if len(m.Projects) < limitP {
		limitP = len(m.Projects)
	}

	// Calculate card width (Total / 3 minus margins)
	pCardWidth := (width / 3) - 4

	for i := 0; i < limitP; i++ {
		p := m.Projects[i]
		name := utils.SafeString(p, "title")
		desc := utils.SafeString(p, "description")

		// Truncate description to keep cards even
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}

		content := fmt.Sprintf("%s\n\n%s", highlight.Render(name), desc)

		card := cardStyle.Copy().
			Width(pCardWidth).
			Height(8).
			Render(content)

		projectCards = append(projectCards, card)
	}
	// Join all project cards horizontally
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, projectCards...) + "\n")

	// --- SECTION 3: SERVICES (3 Cards) ---
	doc.WriteString(sectionTitleStyle.Render("What I Do") + "\n")

	var serviceCards []string
	limitS := 3
	if len(m.Services) < limitS {
		limitS = len(m.Services)
	}

	sCardWidth := (width / 3) - 4

	for i := 0; i < limitS; i++ {
		s := m.Services[i]
		title := utils.SafeString(s, "title")
		price := utils.SafeString(s, "price")

		content := fmt.Sprintf("%s\n\nStarting at %s", highlight.Render(title), subtle.Render(price))

		card := cardStyle.Copy().
			Width(sCardWidth).
			Height(6).
			Render(content)

		serviceCards = append(serviceCards, card)
	}
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, serviceCards...) + "\n")

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
