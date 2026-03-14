package tui

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type CreateUtil struct {
	Cursor   int
	Filename string
	Topic    string
	Input    InputModel
	NoOption bool
	Chosen   bool
}

func CreateUtilInit(filename string, topic string) CreateUtil {
	return CreateUtil{
		Cursor: 0,
		Input:  InitialInputModel(filename, topic),
	}
}

func (m CreateUtil) Init() tea.Cmd {
	return nil
}

func (m CreateUtil) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Chosen {
		model, cmd := m.Input.Update(msg)
		m.Input = model.(InputModel)
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

	if m.Chosen {
		s += m.Input.View().Content
		s += lipgloss.Sprint(lipgloss.NewStyle().PaddingLeft(2).Render(inputHelpMsg()))
	} else {
		// TODO: Need to update this help msg coz we now need to show left and right arrow instead of up and down (For Yes/No selection)
		s += lipgloss.Sprint(lipgloss.NewStyle().PaddingLeft(2).Render(pickerHelpMsg()))
	}

	return tea.NewView(s)
}

type pickerKeyMap struct {
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k pickerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Select, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k pickerKeyMap) FullHelp() [][]key.Binding {
	// No need for this coz I'm not gonna use it... It makes it column wise and its UGLY
	return [][]key.Binding{}
}

var pickerKeys = pickerKeyMap{
	Left: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("←", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("→", "move right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func pickerHelpMsg() string {
	help := help.New()
	helpView := help.ShortHelpView(pickerKeys.ShortHelp())
	return lipgloss.Sprintln(helpView)
}
