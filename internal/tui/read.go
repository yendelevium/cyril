package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
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
		s += fmt.Sprintf("Alias: %s; Path: %s\n", m.Selected.Filename, m.Selected.Filepath)
		fileContent, err := os.ReadFile(m.Selected.Filepath)
		if err != nil {
			s += fmt.Sprintf("Couldn't read file: %v; Error: %v", m.Selected.Filename, err)
			return tea.NewView(s)
		}

		// Get the current terminal width and use that to help deal with border rendering issues in lipgloss
		physicalWidth, _, err := term.GetSize(os.Stdout.Fd())
		if err != nil {
			physicalWidth = 80 // Fallback width just in case
		}
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.RGBColor{
				R: 220,
				G: 155,
				B: 255,
			}).
			Width(physicalWidth - 1).
			PaddingLeft(1).
			PaddingRight(1)

		s += lipgloss.Sprintln(style.Render(fmt.Sprintf("%s", string(fileContent))))
		return tea.NewView(s)
	}

	s += fileIteration(m.Files, m.Cursor)
	s += "  Press q to quit"

	return tea.NewView(s)
}
