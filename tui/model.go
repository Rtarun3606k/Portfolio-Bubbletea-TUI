package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Model struct {
	Width, Height int
	ActiveTab     int // 0: Home, 1: Projects, 2: Experience, 3: Contact
	Loading       bool
	Spinner       spinner.Model

	// Data Storage (Raw BSON maps)
	Projects   []bson.M
	Experience []bson.M
	Services   []bson.M

	//add viewport
	Viewport viewport.Model

	// Contact Form
	EmailInput textinput.Model
	MsgInput   textarea.Model
}

func InitialModel(w, h int) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "your@email.com"
	ti.Focus()

	ta := textarea.New()
	ta.Placeholder = "Your message..."

	return Model{
		Width:      w,
		Height:     h,
		Loading:    true,
		Spinner:    s,
		EmailInput: ti,
		MsgInput:   ta,
	}
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	return InitialModel(pty.Window.Width, pty.Window.Height), []tea.ProgramOption{tea.WithAltScreen()}
}
