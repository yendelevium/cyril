package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type FileData struct {
	Filepath string
	Filename string
}

type ReadModel struct {
	Files    []FileData // items on the to-do list
	Cursor   int        // which to-do list item our Cursor is pointing at
	Selected FileData   // which to-do items are selected
}

func (m ReadModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
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

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m ReadModel) View() tea.View {
	// The header
	var s string

	// IF I've selected a file, display it and quit
	if m.Selected == m.Files[m.Cursor] {
		fileContent, err := os.ReadFile(m.Selected.Filepath)
		if err != nil {
			s += fmt.Sprintf("Couldn't read file: %v; Error: %v", m.Selected.Filename, err)
			return tea.NewView(s)
		}
		s += string(fileContent)
		return tea.NewView(s)
	}

	// Iterate over our Files
	for i, choice := range m.Files {
		Cursor := " "
		if m.Cursor == i {
			Cursor = ">"
		}
		// Render the row
		s += fmt.Sprintf("%s Alias: %s;FilePath: %s;\n", Cursor, choice.Filename, choice.Filepath)
	}

	// The footer
	s += "  Press q to quit"

	// Send the UI for rendering
	return tea.NewView(s)
}
