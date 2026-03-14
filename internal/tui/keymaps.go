package tui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

// Help message footer (currently generalized, might need custom ones later but idk...)
func HelpMsg(keys help.KeyMap) string {
	help := help.New()
	helpView := help.ShortHelpView(keys.ShortHelp())
	return lipgloss.Sprintln(helpView)
}

// generalKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type generalKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k generalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k generalKeyMap) FullHelp() [][]key.Binding {
	// No need for this coz I'm not gonna use it... It makes it column wise and its UGLY
	return [][]key.Binding{}
}

var generalKeys = generalKeyMap{
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
