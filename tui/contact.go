package tui

import (
	"fmt"
	"portfolioTUI/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles matching your web UI
var (
	focusedBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1)
	blurredBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1)
	labelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true).Align(lipgloss.Center)
	subTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Center)
	btnStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("63")).Padding(0, 3).Bold(true)
)

func (m Model) renderContactSection(width int) string {
	// If form was submitted successfully, show a Thank You message
	if m.FormSuccess {
		return lipgloss.Place(width, m.Viewport.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(2).Render(
				lipgloss.JoinVertical(lipgloss.Center,
					titleStyle.Render("Message Sent! 🚀"),
					subTitleStyle.Render("\nThank you for reaching out."),
					lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("\n(Press 'Enter' to send another)"),
				),
			),
		)
	}

	doc := strings.Builder{}

	// --- 1. DYNAMIC WIDTH CALCULATION ---
	// Subtract 6 for padding/borders to ensure no overflow
	fullWidth := width - 6
	if fullWidth < 40 {
		fullWidth = 40
	} // Safety minimum

	// Split width for First/Last name
	halfWidth := (fullWidth - 2) / 2

	// Resize Inputs
	m.FirstNameInput.Width = halfWidth - 3
	m.LastNameInput.Width = halfWidth - 3
	m.EmailInput.Width = fullWidth - 3
	m.MsgInput.SetWidth(fullWidth - 3)

	// --- 2. HEADER ---
	header := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Width(width).Render("Get In Touch"),
		subTitleStyle.Width(width).Render("Have a question or want to work together?"),
	)
	doc.WriteString(header + "\n\n")

	// Helper for border styles
	getStyle := func(i int) lipgloss.Style {
		if m.FocusIndex == i {
			return focusedBorder
		}
		return blurredBorder
	}

	// --- 3. NAME ROW ---
	fName := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("First Name *"),
		getStyle(0).Width(halfWidth).Render(m.FirstNameInput.View()),
	)
	lName := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("Last Name *"),
		getStyle(1).Width(halfWidth).Render(m.LastNameInput.View()),
	)
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, fName, "  ", lName) + "\n\n")

	// --- 4. EMAIL ---
	email := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("Email *"),
		getStyle(2).Width(fullWidth).Render(m.EmailInput.View()),
	)
	doc.WriteString(email + "\n\n")

	// --- 5. USER TYPE ---
	profIcon, studIcon := "( )", "( )"
	if m.UserType == "Professional" {
		profIcon = "(*)"
	} else {
		studIcon = "(*)"
	}

	radioStyle := lipgloss.NewStyle().Padding(1, 0)
	if m.FocusIndex == 3 {
		radioStyle = radioStyle.Foreground(lipgloss.Color("63")).Bold(true)
	}

	doc.WriteString(labelStyle.Render("You are a") + "\n")
	doc.WriteString(radioStyle.Render(fmt.Sprintf("%s Professional    %s Student", profIcon, studIcon)) + "\n\n")

	// --- 6. SERVICES ---
	svcTitle, svcPrice := "-- Select a service --", ""
	if len(m.Services) > 0 && m.SelectedService < len(m.Services) {
		s := m.Services[m.SelectedService]
		svcTitle = utils.SafeString(s, "title")
		svcPrice = utils.SafeString(s, "price")
	}

	// Truncate title
	if len(svcTitle) > fullWidth-20 {
		svcTitle = svcTitle[:fullWidth-20] + "..."
	}
	svcContent := fmt.Sprintf("◄  %-*s  %10s  ►", fullWidth-20, svcTitle, svcPrice)

	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("Select a Service (Optional)"),
		getStyle(4).Width(fullWidth).Render(svcContent),
	) + "\n\n")

	// --- 7. MESSAGE ---
	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("Message *"),
		getStyle(5).Width(fullWidth).Render(m.MsgInput.View()),
	) + "\n\n")

	// --- 8. SUBMIT BUTTON / LOADING ---
	var btnRender string

	if m.ContactLoading {
		// SHOW LOADING SPINNER
		btnRender = lipgloss.JoinHorizontal(lipgloss.Center, m.Spinner.View(), " Sending...")
	} else {
		// SHOW BUTTON
		btnRender = btnStyle.Render("Submit Message ->")
		if m.FocusIndex == 6 {
			btnRender = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("63")).Render(btnRender)
		}
	}

	doc.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, btnRender) + "\n\n")
	return doc.String()
}
