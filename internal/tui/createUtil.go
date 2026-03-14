package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type CreateUtil struct {
	Cursor   int
	Choice   Picker
	Input    InputModel
	NoOption bool
	Chosen   bool
}

type Picker struct {
	Cursor   int
	Filename string
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

var (
	// TODO: CHANGE THESE UGLY STYLES
	focusedStyle = lipgloss.NewStyle().Foreground(BORDERCOLOR)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()

	focusedButton = lipgloss.NewStyle().PaddingLeft(2).Render(focusedStyle.Render("[ Submit ]"))
	blurredButton = lipgloss.Sprintf("%s", lipgloss.NewStyle().PaddingLeft(2).Render(blurredStyle.Render("[ Submit ]")))
)

type InputModel struct {
	FocusIndex int
	Filename   *string
	Topic      *string
	Inputs     []textinput.Model
	CursorMode cursor.Mode
	Quitting   bool
}

func InitialInputModel(filename string, topic string) InputModel {
	m := InputModel{
		Inputs: make([]textinput.Model, 2),
	}

	// TODO: Since only 2 inputs just remove from for loop and write it nomally itself  or what?
	// This works but idk
	for i := range m.Inputs {
		t := textinput.New()
		t.SetWidth(20)
		t.CharLimit = 32

		s := t.Styles()
		s.Cursor.Color = BORDERCOLOR

		// TODO: Change the style for the prompt when focused and blurred
		s.Focused.Prompt = noStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Focused.Text = focusedStyle
		t.SetStyles(s)

		switch i {
		case 0:
			t.Placeholder = "Filename"
			t.SetValue(filename)
			t.Prompt = lipgloss.Sprintf("%s", lipgloss.NewStyle().PaddingLeft(2).Render("Enter the filename: "))
			t.Focus()
		case 1:
			t.Placeholder = "Topic"
			t.SetValue(topic)
			t.Prompt = lipgloss.Sprintf("%s", lipgloss.NewStyle().PaddingLeft(2).Render("Enter the topic: "))
		}

		m.Inputs[i] = t
	}

	return m
}

// TODO: How to translate this into an embedded model?
func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *InputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m InputModel) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	for i, in := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
		if m.CursorMode != cursor.CursorHide && in.Focused() {
			c = in.Cursor()
			if c != nil {
				c.Y += i
			}
		}
	}

	button := &blurredButton
	if m.FocusIndex == len(m.Inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n", *button)

	if m.Quitting {
		b.WriteRune('\n')
	}

	v := tea.NewView(b.String())
	v.Cursor = c
	return v
}

func CreateUtilInit(filename string, topic string) CreateUtil {
	return CreateUtil{
		Cursor: 0,
		Choice: Picker{
			Cursor:   0,
			Filename: filename,
		},
		Input: InitialInputModel(filename, topic),
	}
}

func (m CreateUtil) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m CreateUtil) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Chosen {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.Input.Quitting = true
				return m, tea.Quit

			// Set focus to next input
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if s == "enter" && m.Input.FocusIndex == len(m.Input.Inputs) {
					// TODO: Set like proper values for the filename and shi before quitting here
					return m, tea.Quit
				}

				// Movement
				if s == "up" || s == "shift+tab" {
					if m.Input.FocusIndex > 0 {
						m.Input.FocusIndex--
					}
				} else {
					// On tab, enter and down arrow, move to the next input field
					if m.Input.FocusIndex < len(m.Input.Inputs) {
						m.Input.FocusIndex++
					}
				}

				cmds := make([]tea.Cmd, len(m.Input.Inputs))
				for i := 0; i <= len(m.Input.Inputs)-1; i++ {
					if i == m.Input.FocusIndex {
						// Set focused state
						cmds[i] = m.Input.Inputs[i].Focus()
						continue
					}
					// Remove focused state
					m.Input.Inputs[i].Blur()
				}

				return m, tea.Batch(cmds...)

			// Filenames shouldn't have SPACES
			case "space":
				return m, nil
			}
		}

		// Handle character input and blinking
		cmd := m.Input.updateInputs(msg)
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

type inputKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k inputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k inputKeyMap) FullHelp() [][]key.Binding {
	// No need for this coz I'm not gonna use it... It makes it column wise and its UGLY
	return [][]key.Binding{}
}

var inputKeys = inputKeyMap{
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
		key.WithHelp("↵", "confirm"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
}

func inputHelpMsg() string {
	help := help.New()
	helpView := help.ShortHelpView(inputKeys.ShortHelp())
	return lipgloss.Sprintln(helpView)
}
