package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type InputModel struct {
	FocusIndex int
	Filename   string
	Topic      string
	Inputs     []textinput.Model
	CursorMode cursor.Mode
	Quitting   bool
}

var (
	// TODO: CHANGE THESE UGLY STYLES
	focusedStyle = lipgloss.NewStyle().Foreground(BORDERCOLOR)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()

	focusedButton = lipgloss.NewStyle().PaddingLeft(2).Render(focusedStyle.Render("[ Submit ]"))
	blurredButton = lipgloss.NewStyle().PaddingLeft(2).Render(blurredStyle.Render("[ Submit ]"))
)

func InitialInputModel(filename string, topic string) InputModel {
	m := InputModel{
		Inputs: make([]textinput.Model, 2),
	}

	// TODO: Since only 2 inputs just remove from for loop and write it nomally itself  or what?
	// This works but idk
	for i := range m.Inputs {
		t := textinput.New()
		t.SetWidth(20)
		t.CharLimit = 32

		s := t.Styles()
		s.Cursor.Color = BORDERCOLOR

		// TODO: Change the style for the prompt when focused and blurred
		s.Focused.Prompt = noStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Focused.Text = focusedStyle
		t.SetStyles(s)

		switch i {
		case 0:
			t.Placeholder = "Filename"
			t.SetValue(filename)
			t.Prompt = lipgloss.Sprintf("%s", lipgloss.NewStyle().PaddingLeft(2).Render("Enter the filename: "))
			t.Focus()
		case 1:
			t.Placeholder = "Topic"
			t.SetValue(topic)
			t.Prompt = lipgloss.Sprintf("%s", lipgloss.NewStyle().PaddingLeft(2).Render("Enter the topic: "))
		}

		m.Inputs[i] = t
	}

	return m
}

// Blinking cursor
func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.FocusIndex == len(m.Inputs) {
				// Set the values for the filename and shi before quitting here
				// TODO: INput validation for this?? Can't have empty file/topic names!!
				m.Quitting = true
				return m, tea.Quit
			}

			// Movement
			if s == "up" || s == "shift+tab" {
				if m.FocusIndex > 0 {
					m.FocusIndex--
				}
			} else {
				// On tab, enter and down arrow, move to the next input field
				if m.FocusIndex < len(m.Inputs) {
					m.FocusIndex++
				}
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Focus()
					continue
				}
				// Remove focused state
				m.Inputs[i].Blur()
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

// TODO: INput validation for this?? Can't have empty inputs!
func (m *InputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m InputModel) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	for i, in := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
		if m.CursorMode != cursor.CursorHide && in.Focused() {
			c = in.Cursor()
			if c != nil {
				c.Y += i
			}
		}
	}

	button := &blurredButton
	if m.FocusIndex == len(m.Inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n", *button)

	v := tea.NewView(b.String())
	v.Cursor = c
	return v
}
