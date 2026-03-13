package tui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	// TODO: CHANGE THESE UGLY STYLES
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	quitting   bool
}

func InitialModel() model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	// TODO: Since only 2 inputs just remove from for loop and write it nomally itself  or what?
	// This works but idk
	for i := range m.inputs {
		t := textinput.New()
		t.SetWidth(20)
		t.CharLimit = 32

		s := t.Styles()
		s.Cursor.Color = lipgloss.Color("205")

		// TODO: Change the style for the prompt when focused and blurred
		s.Focused.Prompt = noStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Focused.Text = focusedStyle
		t.SetStyles(s)

		switch i {
		case 0:
			t.Placeholder = "Filename"
			t.SetValue("Filename")
			t.Prompt = "Enter the filename: "
			t.Focus()
		case 1:
			t.Placeholder = "Topic"
			t.SetValue("Topic")
			t.Prompt = "Enter the topic: "
		}

		m.inputs[i] = t
	}

	return m
}

// TODO: How to translate this into an embedded model?
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				// TODO: Set like proper values for the filename and shi before quitting here
				return m, tea.Quit
			}

			// Movement
			if s == "up" || s == "shift+tab" {
				if m.focusIndex > 0 {
					m.focusIndex--
				}
			} else {
				// On tab, enter and down arrow, move to the next input field
				if m.focusIndex < len(m.inputs) {
					m.focusIndex++
				}
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)

		// Filenames shouldn't have SPACES
		case "space":
			return m, nil
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	for i, in := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
		if m.cursorMode != cursor.CursorHide && in.Focused() {
			c = in.Cursor()
			if c != nil {
				c.Y += i
			}
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n", *button)

	// Change this
	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	if m.quitting {
		b.WriteRune('\n')
	}

	v := tea.NewView(b.String())
	v.Cursor = c
	return v
}
