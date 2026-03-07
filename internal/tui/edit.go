package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type EditModel struct {
	Files    []FileData
	Cursor   int
	Selected FileData
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
	return m, nil
}

func (m EditModel) View() tea.View {
	var s string
	if m.NoOption == true {
		s += "No file chosen to edit!\n"
		return tea.NewView(s)
	}

	if m.Selected == m.Files[m.Cursor] {
		s += fmt.Sprintf("[EDITING] Alias: %s;FilePath: %s;\n", m.Selected.Filename, m.Selected.Filepath)
		return tea.NewView(s)
	}
	s += fileIteration(m.Files, m.Cursor)
	s += "  Press q to quit"
	return tea.NewView(s)
}
