package tui

import (
	"fmt"
	"os"
	"strings"

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
