package widgets

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/internal/usecase"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/rivo/tview"
)

type Setting struct {
	Form         *tview.Form
	settingUC    *usecase.SettingUseCase
	errorHandler *service.ErrorHandler
}

func NewSetting(settingUC *usecase.SettingUseCase, errorHandler *service.ErrorHandler) *Setting {
	return &Setting{
		Form: tview.NewForm().
			SetLabelColor(tcell.ColorPurple),
		settingUC:    settingUC,
		errorHandler: errorHandler,
	}
}

func (s *Setting) GenerateInitSetting(tui *service.TUI) {
	setting, err := s.settingUC.Get()
	if err != nil {
		s.errorHandler.ShowErrorWithErr(err, "settingForm")
		return
	}
	s.Form.AddInputField("Show Relative Date(0:Today Only) : ", fmt.Sprint(setting.RelativeDate), 20, nil, nil).
		AddInputField("Person Day : ", fmt.Sprint(setting.PersonDay), 20, nil, nil).
		AddCheckbox("Display As Person Day : ", setting.DisplayAsPersonDay, nil).
		AddInputField("Download Path : ", setting.DownloadPath, 60, nil, nil).
		AddButton("Save", func() {
			s.update()
			s.ReStore(tui)
			tui.SetFocus("menu")
		}).
		AddButton("Cancel", func() {
			s.ReStore(tui)
			tui.SetFocus("menu")
		})
}

func (s *Setting) ReStore(tui *service.TUI) {
	s.Form.Clear(true)
	s.GenerateInitSetting(tui)
}

func (s *Setting) update() {
	relativeDate := s.Form.GetFormItemByLabel("Show Relative Date(0:Today Only) : ").(*tview.InputField).GetText()
	personDay := s.Form.GetFormItemByLabel("Person Day : ").(*tview.InputField).GetText()
	displayAsPersonDay := s.Form.GetFormItemByLabel("Display As Person Day : ").(*tview.Checkbox).IsChecked()
	downloadPath := s.Form.GetFormItemByLabel("Download Path : ").(*tview.InputField).GetText()

	var dateInt, personDayInt int
	var err error
	if dateInt, err = strconv.Atoi(relativeDate); err != nil {
		s.errorHandler.ShowErrorWithErr(err, "settingForm")
		return
	}
	if personDayInt, err = strconv.Atoi(personDay); err != nil {
		s.errorHandler.ShowErrorWithErr(err, "settingForm")
		return
	}

	currentSetting, err := s.settingUC.Get()
	if err != nil {
		s.errorHandler.ShowErrorWithErr(err, "settingForm")
		return
	}
	updatedSetting := &domain.Setting{
		ID:                 currentSetting.ID,
		RelativeDate:       uint(dateInt),
		PersonDay:          uint(personDayInt),
		DisplayAsPersonDay: displayAsPersonDay,
		DownloadPath:       downloadPath,
	}
	if err = s.settingUC.Update(updatedSetting); err != nil {
		s.errorHandler.ShowErrorWithErr(err, "settingForm")
		return
	}
}
