package tui

import (
	"portfolioTUI/config"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Model struct {
	Width, Height int
	ActiveTab     int // 0: Home, ..., 5: Contact
	Loading       bool
	Spinner       spinner.Model

	// Data
	Projects   []bson.M
	Experience []bson.M
	Services   []bson.M
	Blogs      []bson.M

	Viewport viewport.Model

	// --- CONTACT FORM STATE ---
	FirstNameInput textinput.Model
	LastNameInput  textinput.Model
	EmailInput     textinput.Model
	MsgInput       textarea.Model

	UserType        string // "Professional" or "Student"
	SelectedService int    // Index of m.Services
	FocusIndex      int    // 0-6

	// Contact Specific States
	ContactLoading bool // True when "Submit" is clicked
	FormSuccess    bool // True after successful submit
}

func InitialModel(w, h int, data config.AllMessages) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63")) // Purple spinner

	// 1. Initialize Inputs
	fn := textinput.New()
	fn.Placeholder = "Jane"
	fn.CharLimit = 30
	fn.Focus() // Start focused

	ln := textinput.New()
	ln.Placeholder = "Doe"
	ln.CharLimit = 30

	email := textinput.New()
	email.Placeholder = "your.email@example.com"
	email.CharLimit = 50

	ta := textarea.New()
	ta.Placeholder = "Tell me about your project..."
	ta.CharLimit = 500
	ta.SetHeight(5)
	ta.ShowLineNumbers = false

	model := Model{
		Width:   w,
		Height:  h,
		Loading: true,
		Spinner: s,
		// Contact Init
		FirstNameInput: fn,
		LastNameInput:  ln,
		EmailInput:     email,
		MsgInput:       ta,
		UserType:       "Professional",
		FocusIndex:     0,
		ContactLoading: false,
		FormSuccess:    false,
	}

	// 3. LOAD THE PRE-FETCHED DATA IMMEDIATELY
	for _, msg := range data {
		model = updateModelWithData(model, msg)
	}

	// Data is ready, stop loading spinner logic (mostly)
	model.Loading = false

	return model

}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	data := GetOrFetchData()
	pty, _, active := s.Pty()
	if !active {
		return nil, nil
	}
	return InitialModel(pty.Window.Width, pty.Window.Height, data), []tea.ProgramOption{tea.WithAltScreen()}
}
