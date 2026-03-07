package tui

import (
	tea "charm.land/bubbletea/v2"
)

type ReadModel struct {
	Files    []FileData
	Cursor   int
	Selected FileData
}

func (m ReadModel) Init() tea.Cmd {
	return nil
}

func (m ReadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ReadModel) View() tea.View {
	var s string

	// IF I've selected a file, display it and quit
	if m.Selected == m.Files[m.Cursor] {
		s += FileDisplay(m.Selected)
		return tea.NewView(s)
	}

	s += FileIteration(m.Files, m.Cursor)
	s += "  Press q to quit"

	return tea.NewView(s)
}
