package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Dashboard key.Binding
	List      key.Binding
	Prev      key.Binding
	Next      key.Binding
	Quit      key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Dashboard: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "dashboard"),
		),
		List: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "list"),
		),
		Prev: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h/←", "prev month"),
		),
		Next: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l/→", "next month"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}
