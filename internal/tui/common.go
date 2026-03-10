package tui

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

type FileData struct {
	Filepath string
	Filename string
}

// CONSTANTS (but with var coz -> https://stackoverflow.com/a/37984432/22635468)
var BORDERCOLOR lipgloss.RGBColor = lipgloss.RGBColor{
	R: 220,
	G: 155,
	B: 255,
}

func FileIteration(files []FileData, cursor int) string {
	var s string

	// Getting the length for the colour block
	maxLen := 0
	padding := lipgloss.NewStyle().PaddingLeft(2)
	for _, file := range files {
		fileMsg := lipgloss.Sprintln(padding.Render(fmt.Sprintf("Alias: %s; FilePath: %s;", file.Filename, file.Filepath)))
		maxLen = max(maxLen, len(fileMsg))
	}

	// Iterate over our Files
	for i, choice := range files {
		if cursor == i {
			style := padding.
				Background(lipgloss.RGBColor{
					R: 220,
					G: 155,
					B: 255,
				}).
				Width(maxLen)
			s += lipgloss.Sprintln(style.Render(fmt.Sprintf("Alias: %s; FilePath: %s;", choice.Filename, choice.Filepath)))
		} else {
			s += lipgloss.Sprintln(padding.Render(fmt.Sprintf("Alias: %s; FilePath: %s;", choice.Filename, choice.Filepath)))
		}
	}
	return s
}

func FileDisplay(file FileData) string {
	s := fmt.Sprintf("Alias: %s; Path: %s\n", file.Filename, file.Filepath)
	fileContent, err := os.ReadFile(file.Filepath)
	if err != nil {
		s += fmt.Sprintf("Couldn't read file: %v; Error: %v", file.Filename, err)
		return s
	}

	// Get the current terminal width and use that to help deal with border rendering issues in lipgloss
	physicalWidth, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		physicalWidth = 80 // Fallback width just in case
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(BORDERCOLOR).
		Width(physicalWidth).
		PaddingLeft(1).
		PaddingRight(1)

	s += lipgloss.Sprintln(style.Render(fmt.Sprintf("%s", strings.TrimSpace(string(fileContent)))))
	return s

}

// Help message footer (currently generalized, might need custom ones later but idk...)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	// No need for this coz I'm not gonna use it... It makes it column wise and its UGLY
	return [][]key.Binding{}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select item"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func HelpMsg() string {
	help := help.New()
	helpView := help.ShortHelpView(keys.ShortHelp())
	return lipgloss.Sprintln(helpView)
}
