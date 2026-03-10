package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ReadModel struct {
	Files    []FileData
	Cursor   int
	Reply    *FileData
	Selected FileData
	NoOption bool
}

func (m ReadModel) Init() tea.Cmd {
	return nil
}

func (m ReadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
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

			// Reply will be used by the COBRA read command to render the view
			m.Reply.Filename = m.Selected.Filename
			m.Reply.Filepath = m.Selected.Filepath
			return m, tea.Batch(tea.Quit, tea.ClearScreen)
		}
	}
	return m, nil
}

func (m ReadModel) View() tea.View {
	if m.NoOption == true {
		return tea.NewView("No file chosen to read!\n")
	}

	if m.Selected == m.Files[m.Cursor] {
		return tea.NewView("")
	}

	s := FileIteration(m.Files, m.Cursor)

	// Just rendering the help msg according to the file list formatting (padding left 2)
	s += lipgloss.Sprint(lipgloss.NewStyle().PaddingLeft(2).Render(HelpMsg()))
	return tea.NewView(s)
}
