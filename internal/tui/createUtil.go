package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type CreateUtil struct {
	Cursor int
	Reply  *struct {
		Filename string
		Topic    string
	}
	NoOption bool
	Chosen   bool
	Input    InputModel
}

func (m CreateUtil) Init() tea.Cmd {
	return nil
}

func (m CreateUtil) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Chosen {
		model, cmd := m.Input.Update(msg)
		m.Input = model.(InputModel)

		// If I'm quitting, relay the file info
		if m.Input.Quitting {
			m.Reply.Filename = m.Input.Inputs[0].Value()
			m.Reply.Topic = m.Input.Inputs[1].Value()
		}

		return m, cmd
	}

	// IF 'YES' Not chosen
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.NoOption = true
			return m, tea.Quit
		case "left":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "right":
			if m.Cursor < 1 {
				m.Cursor++
			}

		case "enter":
			// IF I haven't chosen any option yet
			if m.Cursor == 1 {
				m.NoOption = true
				return m, tea.Quit
			}
			m.Chosen = true
			return m, nil
		}
	}
	return m, nil
}

func (m CreateUtil) View() tea.View {
	// Just rendering the help msg according to the file list formatting (padding left 2)
	if m.NoOption {
		return tea.NewView("")
	}

	style := lipgloss.NewStyle().PaddingLeft(2)
	s := lipgloss.Sprintln(style.Render(fmt.Sprintf("Couldn't find the file %s!", m.Reply.Filename)))
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

	if m.Chosen {
		s += m.Input.View().Content
		s += lipgloss.Sprint(lipgloss.NewStyle().PaddingLeft(2).Render(HelpMsg(inputKeys)))
	} else {
		// TODO: Need to update this help msg coz we now need to show left and right arrow instead of up and down (For Yes/No selection)
		s += lipgloss.Sprint(lipgloss.NewStyle().PaddingLeft(2).Render(HelpMsg(pickerKeys)))
	}

	return tea.NewView(s)
}
