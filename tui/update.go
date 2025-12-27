package tui

import (
	"log"
	"portfolioTUI/config"
	"portfolioTUI/database"

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

		}

	case DataMsg:
		m = updateModelWithData(m, msg)

	case AllMessages:
		for _, dataMsg := range msg {
			m = updateModelWithData(m, dataMsg)
		}
		m.Loading = false

		// Data is loaded, so we update the viewport with the new data
		// Use 0 width for now if viewport isn't ready, otherwise use existing width
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
		//
		// // 2. Calculate Viewport Height
		// // Screen Height - Header Space (~6 lines) - Margin (2)
		// headerHeight := 6
		// viewportHeight := m.Height - headerHeight
		// if viewportHeight < 5 {
		// 	viewportHeight = 5
		// }

		// 3. Initialize Viewport
		m.Viewport = viewport.New(contentWidth, viewPortHeight)
		m.Viewport.YPosition = headerHeight // Optional, helps with mouse wheel support

		// 4. Set Initial Content
		m.Viewport.SetContent(m.generateConetnt(contentWidth))
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
		m.Experience = msg.Data // assuming 'positions' maps to Experience
	case "services":
		m.Services = msg.Data
	}
	return m
}
