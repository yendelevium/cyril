package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type EditModel struct {
	Files    []FileData // items on the to-do list
	Cursor   int        // which to-do list item our Cursor is pointing at
	Selected FileData   // which to-do items are selected
	Reply    *FileData
	NoOption bool
}

func (m EditModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m EditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.NoOption = true
			return m, tea.Quit

		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down":
			if m.Cursor < len(m.Files)-1 {
				m.Cursor++
			}

		case "enter":
			m.Selected = m.Files[m.Cursor]
			m.Reply.Filename = m.Selected.Filename
			m.Reply.Filepath = m.Selected.Filepath
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m EditModel) View() tea.View {
	// The header
	var s string

	if m.NoOption == true {
		s += "No file chosen to edit!\n"
		return tea.NewView(s)
	}

	if m.Selected == m.Files[m.Cursor] {
		s += fmt.Sprintf("[EDITING] Alias: %s;FilePath: %s;\n", m.Selected.Filename, m.Selected.Filepath)
		return tea.NewView(s)
	}
	maxLen := 0
	for _, file := range m.Files {
		fileMsg := fmt.Sprintf("  Alias: %s;FilePath: %s;\n", file.Filename, file.Filepath)
		maxLen = max(maxLen, len(fileMsg))
	}
	// Iterate over our Files
	for i, choice := range m.Files {
		if m.Cursor == i {
			// TODO: This has a problem where if I press 'q', the colour stays and it looks kindof ugly
			style := lipgloss.NewStyle().Background(lipgloss.RGBColor{
				R: 220,
				G: 155,
				B: 255,
			}).Width(maxLen)
			s += lipgloss.Sprintln(style.Render(fmt.Sprintf("  Alias: %s;FilePath: %s;", choice.Filename, choice.Filepath)))
		} else {
			s += fmt.Sprintf("  Alias: %s;FilePath: %s;\n", choice.Filename, choice.Filepath)
		}
	}

	s += "  Press q to quit"
	return tea.NewView(s)
}
