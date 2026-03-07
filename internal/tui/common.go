package tui

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

type FileData struct {
	Filepath string
	Filename string
}

// TODO: Instead of padding with spaces, actually use lipgloss.Padding()? (pain) (will do later)
func fileIteration(files []FileData, cursor int) string {
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
