package tui

import (
	"portfolioTUI/config"
	"portfolioTUI/database"
	"portfolioTUI/utils"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	// Fetch all 3 collections in parallel
	return tea.Batch(
		m.Spinner.Tick,
		utils.FetchData(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// --- 1. GLOBAL KEY COMMANDS ---
	case tea.KeyMsg:
		// Always allow quitting
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		// --- 2. CONTACT PAGE SPECIFIC LOGIC (Tab 5) ---
		if m.ActiveTab == 5 {

			// A. Handle Success Screen (Reset on Enter)
			if m.FormSuccess {
				if msg.String() == "enter" {
					m.FormSuccess = false
					m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))
				}
				return m, nil
			}

			// B. Handle Form Navigation & Interaction
			// We intercept navigation keys so they don't trigger global tab switching
			switch msg.String() {
			case "up":
				m.FocusIndex--
				if m.FocusIndex < 0 {
					m.FocusIndex = 6
				}
				cmds = append(cmds, m.updateFocus())
				return m, tea.Batch(cmds...)

			case "down", "tab":
				m.FocusIndex++
				if m.FocusIndex > 6 {
					m.FocusIndex = 0
				}
				cmds = append(cmds, m.updateFocus())
				return m, tea.Batch(cmds...)

			// Radio Buttons (Index 3) & Service Select (Index 4) Logic
			case "left", "right":
				if m.FocusIndex == 3 {
					// Toggle User Type
					if m.UserType == "Professional" {
						m.UserType = "Student"
					} else {
						m.UserType = "Professional"
					}
					m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))
					return m, nil

				} else if m.FocusIndex == 4 && len(m.Services) > 0 {
					// Cycle Services
					if msg.String() == "right" {
						m.SelectedService++
						if m.SelectedService >= len(m.Services) {
							m.SelectedService = 0
						}
					} else {
						m.SelectedService--
						if m.SelectedService < 0 {
							m.SelectedService = len(m.Services) - 1
						}
					}
					m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))
					return m, nil
				}
				// If not focusing on Radio/Select, let 'left/right' fall through to Global Navigation

			// Submit Button Logic
			case "enter":
				if m.FocusIndex == 6 && !m.ContactLoading {
					// Validation: Ensure First Name and Email are filled
					if m.FirstNameInput.Value() == "" || m.EmailInput.Value() == "" {
						return m, nil
					}

					// 1. Set Loading State
					m.ContactLoading = true

					// 2. Prepare Data
					fName := m.FirstNameInput.Value()
					lName := m.LastNameInput.Value()
					email := m.EmailInput.Value()
					uType := strings.ToLower(m.UserType)
					msgVal := m.MsgInput.Value()

					// Get Service ID
					svcID := ""
					if len(m.Services) > 0 && m.SelectedService >= 0 && m.SelectedService < len(m.Services) {
						s := m.Services[m.SelectedService]
						svcID = utils.SafeID(s, "_id")
					}

					// 3. Update View immediately to show spinner
					m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))

					// 4. Fire DB Command
					return m, func() tea.Msg {
						time.Sleep(500 * time.Millisecond)
						err := database.InsertContact(fName, lName, email, uType, msgVal, svcID)
						if err != nil {
							return config.FormSubmittedMsg{Success: false}
						}
						return config.FormSubmittedMsg{Success: true}
					}
				}
			}
		}

		// --- 3. GLOBAL TAB NAVIGATION ---
		// Determine if we are typing inside a text box (First, Last, Email, Message)
		isTyping := m.ActiveTab == 5 && (m.FocusIndex == 0 || m.FocusIndex == 1 || m.FocusIndex == 2 || m.FocusIndex == 5)

		// Only allow global navigation if we are NOT typing
		if !isTyping {
			switch msg.String() {

			// FIX: Added "tab" here so you can navigate!
			case "right", "tab":
				if m.ActiveTab < 5 {
					m.ActiveTab++
					m.refreshViewport()
				}
			case "left", "shift+tab":
				if m.ActiveTab > 0 {
					m.ActiveTab--
					m.refreshViewport()
				}

			// Hotkeys
			case "H": // Home
				m.ActiveTab = 0
				m.refreshViewport()
			case "P": // Projects
				m.ActiveTab = 1
				m.refreshViewport()
			case "E": // Experience
				m.ActiveTab = 2
				m.refreshViewport()
			case "S": // Services
				m.ActiveTab = 3
				m.refreshViewport()
			case "B": // Blog
				m.ActiveTab = 4
				m.refreshViewport()
			case "C": // Contact
				m.ActiveTab = 5
				m.refreshViewport()
			}
		}

	// --- 4. FORM SUBMISSION RESULT ---
	case config.FormSubmittedMsg:
		m.ContactLoading = false
		if msg.Success {
			m.FormSuccess = true
			m.FirstNameInput.SetValue("")
			m.LastNameInput.SetValue("")
			m.EmailInput.SetValue("")
			m.MsgInput.SetValue("")
			m.FocusIndex = 0
		}
		m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))
		return m, nil

	// --- 5. DATA FETCHING ---
	case config.DataMsg:
		m = updateModelWithData(m, msg)
		newCmds := utils.GenerateImagesCmds(msg.Type, msg.Data)
		cmds = append(cmds, newCmds...)

	case config.AllMessages:
		for _, dataMsg := range msg {
			m = updateModelWithData(m, dataMsg)
			newCmds := utils.GenerateImagesCmds(dataMsg.Type, dataMsg.Data)
			cmds = append(cmds, newCmds...)
		}
		m.Loading = false
		m.refreshViewport()

	// --- 6. WINDOW RESIZE ---
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		headerHeight := 5
		footerHeight := 5
		verticalMargin := headerHeight + footerHeight
		viewPortHeight := m.Height - verticalMargin
		if viewPortHeight < 5 {
			viewPortHeight = 5
		}

		contentWidth := int(float64(m.Width) * 0.8)
		if contentWidth < 40 {
			contentWidth = 40
		}

		m.Viewport = viewport.New(contentWidth, viewPortHeight)
		m.Viewport.YPosition = headerHeight
		m.Viewport.SetContent(m.generateConetnt(contentWidth))

	// --- 7. IMAGE GENERATION RESULT ---
	case utils.AsciiIamge:
		var shouldRefresh bool
		switch msg.CollectionName {
		case "projects":
			if msg.Index < len(m.Projects) {
				m.Projects[msg.Index]["ascii_art"] = msg.Art
				if m.ActiveTab == 0 || m.ActiveTab == 1 {
					shouldRefresh = true
				}
			}
		case "positions":
			if msg.Index < len(m.Experience) {
				m.Experience[msg.Index]["ascii_art"] = msg.Art
				if m.ActiveTab == 2 {
					shouldRefresh = true
				}
			}
		case "blogs":
			if msg.Index < len(m.Blogs) {
				m.Blogs[msg.Index]["ascii_art"] = msg.Art
				if m.ActiveTab == 0 || m.ActiveTab == 4 {
					shouldRefresh = true
				}
			}
		}
		if shouldRefresh {
			m.refreshViewport()
		}
	}

	// --- 8. UPDATE BUBBLES ---
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	// Inputs (Only update if on Contact Page)
	if m.ActiveTab == 5 {
		m.FirstNameInput, cmd = m.FirstNameInput.Update(msg)
		cmds = append(cmds, cmd)
		m.LastNameInput, cmd = m.LastNameInput.Update(msg)
		cmds = append(cmds, cmd)
		m.EmailInput, cmd = m.EmailInput.Update(msg)
		cmds = append(cmds, cmd)
		m.MsgInput, cmd = m.MsgInput.Update(msg)
		cmds = append(cmds, cmd)

		// Live View Update: Optimized to only render when actually typing
		if keyMsg, ok := msg.(tea.KeyMsg); ok && !m.ContactLoading && isTypingInput(keyMsg) {
			m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))
		}
	}

	return m, tea.Batch(cmds...)
}

// --- HELPER FUNCTIONS ---

// updateFocus handles blurring/focusing inputs based on m.FocusIndex
func (m *Model) updateFocus() tea.Cmd {
	// 1. Blur all
	m.FirstNameInput.Blur()
	m.LastNameInput.Blur()
	m.EmailInput.Blur()
	m.MsgInput.Blur()

	var cmd tea.Cmd

	// 2. Focus specific
	switch m.FocusIndex {
	case 0:
		cmd = m.FirstNameInput.Focus()
	case 1:
		cmd = m.LastNameInput.Focus()
	case 2:
		cmd = m.EmailInput.Focus()
	case 5:
		cmd = m.MsgInput.Focus()
	}

	// 3. Refresh view to update border colors
	m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))

	return cmd
}

// refreshViewport regenerates the current tab's content
func (m *Model) refreshViewport() {
	// If on Contact page (5), use the special render function
	// Otherwise use the generic generator
	if m.ActiveTab == 5 {
		m.Viewport.SetContent(m.renderContactSection(m.Viewport.Width))
	} else {
		// Note: Keeping your original typo 'generateConetnt' to ensure compatibility
		m.Viewport.SetContent(m.generateConetnt(m.Viewport.Width))
	}
	m.Viewport.GotoTop()
}

func isTypingInput(msg tea.KeyMsg) bool {
	// Simple check if the key is likely a printable character
	return len(msg.String()) == 1 || msg.String() == "backspace" || msg.String() == "space"
}

// Helper function to keep the switch clean
func updateModelWithData(m Model, msg config.DataMsg) Model {
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
