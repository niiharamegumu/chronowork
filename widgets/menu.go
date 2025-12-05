package widgets

import (
	"github.com/niiharamegumu/chronowork/internal/usecase"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/niiharamegumu/chronowork/util/timeutil"
	"github.com/rivo/tview"
)

type Menu struct {
	List      *tview.List
	settingUC *usecase.SettingUseCase
}

func NewMenu(settingUC *usecase.SettingUseCase) *Menu {
	return &Menu{
		List:      tview.NewList(),
		settingUC: settingUC,
	}
}

func (m *Menu) addListItem(text string, shortcut rune, selected func()) *Menu {
	m.List.AddItem(text, "", shortcut, selected)
	return m
}

func (m *Menu) GenerateInitMenu(tui *service.TUI, work *Work, setting *Setting, project *Project) *Menu {
	m.addListItem("Works", 'w', func() {
		relativeDays := m.getRelativeDays()
		work.ReStoreTable(timeutil.RelativeStartTimeWithDays(relativeDays), timeutil.TodayEndTime())
		tui.ChangeToPage("work")
		tui.SetFocus("mainWorkContent")
	})
	m.addListItem("Projects", 'p', func() {
		project.RestoreTable()
		tui.ChangeToPage("project")
		tui.SetFocus("projectTable")
	})
	m.addListItem("Tags", 't', func() {
		tui.ChangeToPage("tag")
		tui.SetFocus("tagTable")
	})
	m.addListItem("Export", 'e', func() {
		tui.ChangeToPage("export")
		tui.SetFocus("exportForm")
	})
	m.addListItem("Setting", 's', func() {
		setting.ReStore(tui)
		tui.ChangeToPage("setting")
		tui.SetFocus("settingForm")
	})
	m.addListItem("Quit", 'q', tui.Quit)
	return m
}

func (m *Menu) getRelativeDays() int {
	setting, err := m.settingUC.Get()
	if err != nil {
		return 0
	}
	return int(setting.RelativeDate)
}
