package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
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

	// TODO: Instead of padding with spaces, actually use lipgloss.Padding()? (pain) (will do later)
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

	// The footer
	s += "  Press q to quit"

	// Send the UI for rendering
	return tea.NewView(s)
}
