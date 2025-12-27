package tui

import (
	"log"
	"portfolioTUI/config"
	"portfolioTUI/database"
	"portfolioTUI/utils"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Messages
type DataMsg struct {
	Type string // "projects", "experience", "services"
	Data []bson.M
}

// list of all messages
type AllMessages []DataMsg

func fetchData() tea.Cmd {
	return func() tea.Msg {
		var alldata AllMessages
		// Use your new generic function
		for _, collections := range config.Collection { // Ensure config.Collections matches your config file
			data, err := database.GetALLFromCollection(collections, collections)
			if err != nil {
				log.Println("Error fetching", collections, err)
				continue
			}
			log.Println("Fetched", len(data), "items from", collections)
			alldata = append(alldata, DataMsg{Type: collections, Data: data})
		}
		return alldata
	}
}

func (m Model) Init() tea.Cmd {
	// Fetch all 3 collections in parallel
	return tea.Batch(
		m.Spinner.Tick,
		fetchData(),
	)
}

func generateImagesCmds(dataType string, data []bson.M) []tea.Cmd {
	var cmds []tea.Cmd
	defaultImage := config.DEFAULTIMAGEURL

	for i, item := range data {
		var url string
		var key string

		// 1. Determine the correct key for this collection
		switch dataType {
		case "projects":
			key = "imageUrl"
		case "positions":
			key = "logoUrl" // Ensure this matches your DB
		case "blogs":
			key = "featuredImage" // Ensure this matches your DB
		default:
			continue
		}

		// 2. Safe Get: Get the string, or "" if missing
		if val, ok := item[key].(string); ok {
			url = val
		}

		// 3. Fallback Logic (The Fix)
		// If URL is empty OR too short, use the default
		if len(url) < 5 {
			url = defaultImage
		}

		// 4. Generate the Command
		// Now we always have a valid URL (either original or default)
		cmds = append(cmds, utils.GenerateAsciiImage(url, dataType, i, 30, 15))
	}
	return cmds
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// NAVIGATION LOGIC
		// When we switch tabs, we must REFRESH the Viewport content and reset scroll to top
		case "tab", "right":
			if m.ActiveTab < 3 {
				m.ActiveTab++
				content := m.generateConetnt(m.Viewport.Width)
				m.Viewport.SetContent(content)
				m.Viewport.GotoTop()
			}
		case "shift+tab", "left":
			if m.ActiveTab > 0 {
				m.ActiveTab--
				content := m.generateConetnt(m.Viewport.Width)
				m.Viewport.SetContent(content)
				m.Viewport.GotoTop()
			}

		case "H":
			m.ActiveTab = 0
			content := m.generateConetnt(m.Viewport.Width)
			m.Viewport.SetContent(content)
			m.Viewport.GotoTop()
		case "P":
			m.ActiveTab = 1
			content := m.generateConetnt(m.Viewport.Width)
			m.Viewport.SetContent(content)
			m.Viewport.GotoTop()
		case "E":
			m.ActiveTab = 2
			content := m.generateConetnt(m.Viewport.Width)
			m.Viewport.SetContent(content)
			m.Viewport.GotoTop()
		case "C":
			m.ActiveTab = 3
			content := m.generateConetnt(m.Viewport.Width)
			m.Viewport.SetContent(content)
			m.Viewport.GotoTop()

			// Refresh content after any tab change
			if msg.Type == tea.KeyRunes || msg.Type == tea.KeySpace { // Only refresh if tab actually changed logic ran
				// (Simplified: Just ensure content is refreshed if ActiveTab changed)
				m.Viewport.SetContent(m.generateConetnt(m.Viewport.Width))
				m.Viewport.GotoTop()
			}
		}

	case DataMsg:
		m = updateModelWithData(m, msg)

		// Use the new helper to trigger images for this specific data
		newCmds := generateImagesCmds(msg.Type, msg.Data)
		cmds = append(cmds, newCmds...)

	case AllMessages:
		for _, dataMsg := range msg {
			m = updateModelWithData(m, dataMsg)
			newCmds := generateImagesCmds(dataMsg.Type, dataMsg.Data)
			cmds = append(cmds, newCmds...)
		}
		m.Loading = false

		// Data is loaded, so we update the viewport with the new data
		width := m.Viewport.Width
		if width == 0 {
			width = 80 // Fallback
		}
		m.Viewport.SetContent(m.generateConetnt(width))

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerHeight := 5
		footerHeight := 5
		verticalMargin := headerHeight + footerHeight

		viewPortHeight := m.Height - verticalMargin // 1. Calculate Content Size (80% width)

		if viewPortHeight < 5 {
			viewPortHeight = 5

		}

		// This matches your View logic so the text wraps correctly
		contentWidth := int(float64(m.Width) * 0.8)
		if contentWidth < 40 {
			contentWidth = 40
		}

		// 3. Initialize Viewport
		m.Viewport = viewport.New(contentWidth, viewPortHeight)
		m.Viewport.YPosition = headerHeight // Optional, helps with mouse wheel support

		// 4. Set Initial Content
		m.Viewport.SetContent(m.generateConetnt(contentWidth))

	case utils.AsciiIamge:
		switch msg.CollectionName {

		case "projects":
			if msg.Index < len(m.Projects) {
				m.Projects[msg.Index]["ascii_art"] = msg.Art
			}

		case "positions":
			if msg.Index < len(m.Experience) {
				m.Experience[msg.Index]["ascii_art"] = msg.Art
			}

		case "blogs":
			if msg.Index < len(m.Blogs) {
				m.Blogs[msg.Index]["ascii_art"] = msg.Art
			}
		}
	}

	// HANDLE SCROLLING
	// Pass messages to viewport (handles j, k, up, down, pgup, pgdown, mousewheel)
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	// HANDLE SPINNER
	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// Helper function to keep the switch clean
func updateModelWithData(m Model, msg DataMsg) Model {
	switch msg.Type {
	case "projects":
		m.Projects = msg.Data
	case "positions":
		m.Experience = msg.Data
	case "services":
		m.Services = msg.Data
	case "blogs":
		m.Blogs = msg.Data
	}
	return m
}
