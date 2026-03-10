package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

type FileSelector struct {
	Files    []FileData
	Cursor   int
	MaxWidth int // The width of ther terminal /4; Also doubles as the max line width
	Height   int // Height of the text block
}

func (m FileSelector) Init() tea.Cmd {
	return nil
}

func (m FileSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ??
	return m, nil
}

// TODO: Make this view SCROLLABLE instead to not go down in the terminal using the physicalHeight thing?
// Like set a maxSize based on terminal max_height (as its measured by no.of of lines in the window (including older shell prompts)) or an arbitrarty value like 10?
// TODO: Print the path to the file (maybe just topicname and the filename) to help know which one is an alias and which isnt?
func (m FileSelector) View() tea.View {
	msgLen := m.MaxWidth/4 - 2
	msgs := []string{}
	for _, file := range m.Files {
		fileMsg := fmt.Sprintf(" %s ", file.Filename)
		if len(fileMsg) > msgLen {
			fileMsg = fileMsg[:msgLen-2] + ".."
		}
		msgs = append(msgs, fileMsg)
	}

	var s string
	for i, choice := range msgs {
		if m.Cursor == i {
			// TODO: This has a problem where if I press 'q', the colour stays and it looks kindof ugly
			style := lipgloss.NewStyle().Background(lipgloss.RGBColor{
				R: 220,
				G: 155,
				B: 255,
			}).Width(msgLen)
			s += lipgloss.Sprintln(style.Render(choice))
		} else {
			s += lipgloss.Sprintln(choice)
		}
	}

	// WAIT Idk is this it??
	return tea.NewView(lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).Render(s))
}

type FileDisplayer struct {
	FileContents []string
	Cursor       int
	MaxWidth     int // Need this to be terminal width*3/4
	MaxHeight    int // This should be the max height of the textblock WITH the borders. I don't want any border cutoff business either
}

func (m FileDisplayer) Init() tea.Cmd {
	return nil
}

func (m FileDisplayer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ??
	return m, nil
}

// TODO: Make this view SCROLLABLE instead to not go down in the terminal using the physicalHeight thing
func (m FileDisplayer) View() tea.View {
	// WAIT Idk is this it??
	padding := lipgloss.NewStyle().PaddingLeft(1)
	border := lipgloss.NewStyle().
		// Just so you know
		BorderStyle(lipgloss.NormalBorder()).
		PaddingLeft(1).
		PaddingRight(1).
		BorderForeground(BORDERCOLOR).
		Width(m.MaxWidth).
		Height(m.MaxHeight).
		MaxHeight(m.MaxHeight)

	return tea.NewView(padding.Render(border.Render(m.FileContents[m.Cursor])))
}

type ListModel struct {
	Selector   FileSelector
	Displayer  FileDisplayer
	TotalFiles int
	Cursor     int
	Selected   FileData
}

func ListModelInitialize(files []FileData, fileContents []string) ListModel {

	// Width control + display truncation
	physicalWidth, physicalHeight, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		physicalWidth = 80  // Fallback width just in case
		physicalHeight = 80 // Fallback height just in case
	}

	// Truncating the content displayed based on terminal width and height!
	// The main problem is with height. The \n's are confusing me here.
	// The max length of a line, accounting for borders and all is the Width - 4
	// Max height is m.MaxHeight - 2
	// TODO: Also no because there has to be a better way to do this?? What in the DSA is this TvT
	// Also doing this here because I don't wanna re calc everytime I switch between files (if I put this in FileDisplayer.View())

	// TODO: Better maxwidth calculations coz it has some leftoverspace on the side and its annoying me a bit
	lineWidth := (3*physicalWidth-2)/4 - 1
	lineHeight := len(files)

	// Girl I don't even know what I wrote but it works (ok its not that bad but still)
	for idx, content := range fileContents {
		charCount := 0
		lineCount := 0
		for slice, char := range content {
			if char == '\n' {
				lineCount++
				charCount = 0
			} else {
				charCount++
			}

			if charCount == lineWidth {
				charCount = 0
				lineCount++
			}
			if lineCount == lineHeight {
				// Truncate
				content = content[:slice-3] + "..."
				break
			}
		}
		fileContents[idx] = content
	}

	// Make the model
	return ListModel{
		Selector: FileSelector{
			Files:    files,
			Cursor:   0,
			MaxWidth: physicalWidth,
			Height:   physicalHeight,
		},
		Displayer: FileDisplayer{
			FileContents: fileContents,
			Cursor:       0,
			MaxWidth:     lineWidth,
			MaxHeight:    len(files) + 2, // +2 to take account for the padding in the Selector view
		},
		Cursor:     0,
		TotalFiles: len(files),
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
			m.Displayer.Cursor = m.Cursor
			m.Selector.Cursor = m.Cursor
		case "down":
			if m.Cursor < m.TotalFiles-1 {
				m.Cursor++
			}
			m.Displayer.Cursor = m.Cursor
			m.Selector.Cursor = m.Cursor
		}
	}
	return m, nil
}

func (m ListModel) View() tea.View {
	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Left, m.Selector.View().Content, m.Displayer.View().Content))
}
