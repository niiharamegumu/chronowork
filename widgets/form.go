package widgets

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/internal/usecase"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/niiharamegumu/chronowork/util/timeutil"
	"github.com/rivo/tview"
)

var notSelectText = "Not Select"

type Form struct {
	Form          *tview.Form
	chronoWorkUC  *usecase.ChronoWorkUseCase
	projectTypeUC *usecase.ProjectTypeUseCase
	errorHandler  *service.ErrorHandler
}

func NewForm(chronoWorkUC *usecase.ChronoWorkUseCase, projectTypeUC *usecase.ProjectTypeUseCase, errorHandler *service.ErrorHandler) *Form {
	form := &Form{
		Form: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple).
			SetFieldTextColor(tcell.ColorGray).
			SetFieldBackgroundColor(tcell.ColorWhite),
		chronoWorkUC:  chronoWorkUC,
		projectTypeUC: projectTypeUC,
		errorHandler:  errorHandler,
	}
	return form
}

func (f *Form) GenerateInitForm(tui *service.TUI, work *Work, relativeDays int) *Form {
	f.ConfigureStoreForm(tui, work, relativeDays)
	return f
}

func (f *Form) FormCapture(tui *service.TUI) {
	f.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlB:
			tui.SetFocus("mainWorkContent")
		}
		return event
	})
}

func (f *Form) ResetForm() {
	f.Form.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
	f.Form.GetFormItemByLabel("Project").(*tview.DropDown).SetCurrentOption(0)
	f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).SetOptions([]string{notSelectText}, nil).SetCurrentOption(0)
}

func (f *Form) ConfigureStoreForm(tui *service.TUI, work *Work, relativeDays int) {
	f.Form.
		AddInputField("Title", "", 50, nil, nil).
		AddDropDown("Project", append([]string{notSelectText}, f.projectTypeUC.GetAllNames()...), 0, f.projectDropDownChanged).
		AddDropDown("Tags", []string{notSelectText}, 0, nil).
		AddButton("Store", func() {
			if err := f.store(); err != nil {
				f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
				return
			}
			if err := work.ReStoreTable(timeutil.RelativeStartTimeWithDays(relativeDays), timeutil.TodayEndTime()); err != nil {
				f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
				return
			}
			tui.SetFocus("mainWorkContent")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("mainWorkContent")
		})
}

func (f *Form) configureUpdateForm(tui *service.TUI, work *Work, chronoWork *domain.ChronoWork, relativeDays int) {
	projectOptions := append([]string{notSelectText}, f.projectTypeUC.GetAllNames()...)
	tagsOptions := []string{notSelectText}
	f.Form.AddInputField("Title", chronoWork.Title, 50, nil, nil).
		AddDropDown("Project", projectOptions, 0, f.projectDropDownChanged).
		AddDropDown("Tags", tagsOptions, 0, nil)

	if chronoWork.ProjectType != nil && chronoWork.ProjectType.Name != "" {
		projectType, err := f.projectTypeUC.FindByName(chronoWork.ProjectType.Name)
		if err != nil {
			f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
			return
		}
		for i, projectOption := range projectOptions {
			if projectOption == chronoWork.ProjectType.Name {
				f.Form.GetFormItemByLabel("Project").(*tview.DropDown).SetCurrentOption(i)
				break
			}
		}

		tagsOptions = append(tagsOptions, projectType.GetTagNames()...)
		f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).SetOptions(tagsOptions, nil)
		if chronoWork.Tag != nil && chronoWork.Tag.Name != "" {
			for i, tagOption := range tagsOptions {
				if tagOption == chronoWork.Tag.Name {
					f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).SetCurrentOption(i)
					break
				}
			}
		}
	}

	f.Form.AddButton("Update", func() {
		if err := f.update(chronoWork); err != nil {
			f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
			return
		}
		if err := work.ReStoreTable(timeutil.RelativeStartTimeWithDays(relativeDays), timeutil.TodayEndTime()); err != nil {
			f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
			return
		}
		tui.SetFocus("mainWorkContent")
	}).
		AddButton("Cancel", func() {
			tui.SetFocus("mainWorkContent")
		})
}

func (f *Form) configureTimerForm(tui *service.TUI, work *Work, chronoWork *domain.ChronoWork, relativeDays int) {
	hour := chronoWork.TotalSeconds / 3600
	minute := (chronoWork.TotalSeconds - hour*3600) / 60
	second := chronoWork.TotalSeconds - hour*3600 - minute*60

	f.Form.AddInputField("Hour(0-)", fmt.Sprint(hour), 20, nil, nil).
		AddInputField("Minute(0-59)", fmt.Sprint(minute), 20, nil, nil).
		AddInputField("Second(0-59)", fmt.Sprint(second), 20, nil, nil).
		AddButton("Reset", func() {
			if err := f.resetTimer(chronoWork); err != nil {
				f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
				return
			}
			if err := work.ReStoreTable(timeutil.RelativeStartTimeWithDays(relativeDays), timeutil.TodayEndTime()); err != nil {
				f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
				return
			}
			tui.SetFocus("mainWorkContent")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("mainWorkContent")
		})
}

func (f *Form) projectDropDownChanged(option string, optionIndex int) {
	var tagsOptions []string
	if f.Form.GetFormItemByLabel("Tags") == nil {
		return
	}
	projectType, err := f.projectTypeUC.FindByName(option)
	if err != nil {
		tagsOptions = []string{notSelectText}
	} else {
		tagsOptions = append([]string{notSelectText}, projectType.GetTagNames()...)
	}
	f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).
		SetOptions(tagsOptions, nil).
		SetCurrentOption(0)
}

func (f *Form) store() error {
	title := f.Form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
	_, projectVal := f.Form.GetFormItemByLabel("Project").(*tview.DropDown).GetCurrentOption()
	_, tagVal := f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).GetCurrentOption()

	if title == "" {
		return nil
	}

	var projectTypeID uint
	var tagID uint
	if projectVal != notSelectText {
		projectType, err := f.projectTypeUC.FindByName(projectVal)
		if err != nil {
			f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
			return err
		}
		projectTypeID = projectType.ID
		if tagVal != notSelectText {
			for _, tag := range projectType.Tags {
				if tag.Name == tagVal {
					tagID = tag.ID
				}
			}
		}
	}

	// 4. ユースケースを呼び出す（ビジネスロジックに委譲）
	if _, err := f.chronoWorkUC.Create(title, projectTypeID, tagID); err != nil {
		f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
		return err
	}
	return nil
}

func (f *Form) update(chronoWork *domain.ChronoWork) error {
	title := f.Form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
	_, projectVal := f.Form.GetFormItemByLabel("Project").(*tview.DropDown).GetCurrentOption()
	_, tagVal := f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).GetCurrentOption()

	if title == "" {
		return nil
	}
	var projectTypeID uint = 0
	var tagID uint = 0
	if projectVal != notSelectText {
		projectType, err := f.projectTypeUC.FindByName(projectVal)
		if err != nil {
			f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
			return err
		}
		projectTypeID = projectType.ID
		if tagVal != notSelectText {
			for _, tag := range projectType.Tags {
				if tag.Name == tagVal {
					tagID = tag.ID
				}
			}
		}
	}
	if err := f.chronoWorkUC.Update(chronoWork.ID, title, projectTypeID, tagID); err != nil {
		f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
		return err
	}

	return nil
}

func (f *Form) resetTimer(chronoWork *domain.ChronoWork) error {
	hour := f.Form.GetFormItemByLabel("Hour(0-)").(*tview.InputField).GetText()
	minute := f.Form.GetFormItemByLabel("Minute(0-59)").(*tview.InputField).GetText()
	second := f.Form.GetFormItemByLabel("Second(0-59)").(*tview.InputField).GetText()

	var hourInt, minuteInt, secondInt uint64
	var err error
	if hourInt, err = strconv.ParseUint(hour, 10, 16); err != nil {
		return err
	}
	if minuteInt, err = strconv.ParseUint(minute, 10, 16); err != nil {
		return err
	}
	if secondInt, err = strconv.ParseUint(second, 10, 16); err != nil {
		return err
	}
	totalSeconds := int(hourInt*3600 + minuteInt*60 + secondInt)
	err = f.chronoWorkUC.UpdateTotalSeconds(chronoWork.ID, totalSeconds)
	if err != nil {
		f.errorHandler.ShowErrorWithErr(err, "mainWorkForm")
		return err
	}

	return nil
}
