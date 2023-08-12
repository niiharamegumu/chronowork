package widgets

import (
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
)

type Menu struct {
	List *tview.List
}

func NewMenu() *Menu {
	return &Menu{
		List: tview.NewList(),
	}
}

func (m *Menu) AddListItem(text string, shortcut rune, selected func()) *Menu {
	m.List.AddItem(text, "", shortcut, selected)
	return m
}

func GenerateInitMenu(tui *service.TUI) *Menu {
	menu := NewMenu()
	menu.AddListItem("Works", 'w', func() {
		tui.SetFocus("mainContent")
	})
	menu.AddListItem("Quit", 'q', tui.Quit)
	return menu
}