package tui

import (
	"fmt"
	"os"

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

// TODO: Instead of padding with spaces, actually use lipgloss.Padding()? (pain) (will do later)
func FileIteration(files []FileData, cursor int) string {
	var s string
	maxLen := 0
	for _, file := range files {
		fileMsg := fmt.Sprintf("  Alias: %s;FilePath: %s;\n", file.Filename, file.Filepath)
		maxLen = max(maxLen, len(fileMsg))
	}
	// Iterate over our Files
	for i, choice := range files {
		if cursor == i {
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
	return s
}

func FileDisplay(file FileData) string {
	var s string
	s += fmt.Sprintf("Alias: %s; Path: %s\n", file.Filename, file.Filepath)
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
		Width(physicalWidth - 1).
		PaddingLeft(1).
		PaddingRight(1)

	s += lipgloss.Sprintln(style.Render(fmt.Sprintf("%s", string(fileContent))))
	return s

}
