package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type CreateUtil struct {
	Cursor   int
	Choice   Picker
	Create   Creator
	NoOption bool
	Chosen   bool
}

type Picker struct {
	Cursor   int
	Filename string
}

type Creator struct {
	Cursor   int
	Filename string
	Topic    string
}

func (m Picker) View() string {
	style := lipgloss.NewStyle().PaddingLeft(2)
	s := lipgloss.Sprintln(style.Render(fmt.Sprintf("Couldn't find the file %s!", m.Filename)))
	prompt := lipgloss.Sprintln(style.Render("Would you like to create the file?"))

	activeButtonStyle := lipgloss.NewStyle().Background(BORDERCOLOR).PaddingLeft(1).PaddingRight(1)
	buttonStyle := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
	yesButton := buttonStyle.Render("Yes")
	noButton := buttonStyle.Render("No")

	// Render the active button
	if m.Cursor == 0 {
		yesButton = activeButtonStyle.Render("Yes")
	} else {
		noButton = activeButtonStyle.Render("No")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, lipgloss.NewStyle().PaddingLeft(1).Render(yesButton), lipgloss.NewStyle().PaddingLeft(1).Render(noButton))
	s += lipgloss.Sprintln(lipgloss.JoinHorizontal(lipgloss.Left, prompt, buttons))

	return s
}

func CreateUtilInit(filename string, topic string) CreateUtil {
	return CreateUtil{
		Cursor: 0,
		Choice: Picker{
			Cursor:   0,
			Filename: filename,
		},
		Create: Creator{
			Filename: filename,
			Topic:    topic,
			Cursor:   0,
		},
	}
}

func (m CreateUtil) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m CreateUtil) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.NoOption = true
			return m, tea.Quit
		case "left":
			if !m.Chosen {
				if m.Cursor > 0 {
					m.Cursor--
					m.Choice.Cursor = m.Cursor
				}
			}

		case "right":
			if !m.Chosen {
				if m.Cursor < 1 {
					m.Cursor++
					m.Choice.Cursor = m.Cursor
				}
			}
		case "enter":

			// IF I haven't chosen any option yet
			if !m.Chosen {
				if m.Cursor == 1 {
					m.NoOption = true
					return m, tea.Quit
				}
				m.Chosen = true
				return m, nil
			}

		}
	}
	return m, nil
}

func (m CreateUtil) View() tea.View {
	// Just rendering the help msg according to the file list formatting (padding left 2)
	if m.NoOption {
		return tea.NewView("")
	}
	s := m.Choice.View()

	if m.Chosen {
		s += lipgloss.Sprint(lipgloss.NewStyle().PaddingLeft(2).Render(HelpMsg()))
	} else {
		// TODO: Need to update this help msg coz we now need to show left and right arrow instead of up and down (For Yes/No selection)
		s += lipgloss.Sprint(lipgloss.NewStyle().PaddingLeft(2).Render(HelpMsg()))
	}

	return tea.NewView(s)
}
