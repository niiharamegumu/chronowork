package widgets

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/internal/usecase"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/niiharamegumu/chronowork/util/strutil"
	"github.com/rivo/tview"
)

var (
	projectHeader = []string{
		"ID",
		"Name",
		"Tags",
	}
)

type Project struct {
	Layout        *tview.Grid
	Form          *tview.Form
	ReadOnlyForm  *tview.Form
	Table         *tview.Table
	projectTypeUC *usecase.ProjectTypeUseCase
	tagUC         *usecase.TagUseCase
	chronoWorkUC  *usecase.ChronoWorkUseCase
	errorHandler  *service.ErrorHandler
}

func NewProject(projectTypeUC *usecase.ProjectTypeUseCase, tagUC *usecase.TagUseCase, chronoWorkUC *usecase.ChronoWorkUseCase, errorHandler *service.ErrorHandler) *Project {
	return &Project{
		Layout: tview.NewGrid().
			SetRows(0, 0).
			SetColumns(0, 0).
			SetBorders(true),
		Form: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple).
			SetFieldTextColor(tcell.ColorGray).
			SetFieldBackgroundColor(tcell.ColorWhite),
		ReadOnlyForm: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple).
			SetFieldTextColor(tcell.ColorGray).
			SetFieldBackgroundColor(tcell.ColorWhite),
		Table: tview.NewTable().
			SetSelectable(true, false).
			SetFixed(1, 1),
		projectTypeUC: projectTypeUC,
		tagUC:         tagUC,
		chronoWorkUC:  chronoWorkUC,
		errorHandler:  errorHandler,
	}
}

func (p *Project) GenerateInitProject(tui *service.TUI) *Project {
	p.setStoreProjectForm(tui)
	p.RestoreTable()

	p.Layout.AddItem(p.Form, 0, 0, 1, 1, 0, 0, false)
	p.Layout.AddItem(p.ReadOnlyForm, 0, 1, 1, 1, 0, 0, false)
	p.Layout.AddItem(p.Table, 1, 0, 1, 2, 0, 0, true)

	p.tableCapture(tui)
	p.formCapture(tui)
	return p
}

func (p *Project) RestoreTable() {
	p.Table.Clear()
	p.setTable()
}

func (p *Project) setStoreProjectForm(tui *service.TUI) {
	p.Form.Clear(true)
	p.ReadOnlyForm.Clear(true)

	p.ReadOnlyForm.AddTextArea("Selected Tags", "", 50, 5, 0, nil)
	tags := p.tagUC.GetAllNames()
	tags = append([]string{notSelectText}, tags...)
	p.Form.AddInputField("Project Name : ", "", 50, nil, nil).
		AddDropDown("Tags : ", tags, 0, func(option string, optionIndex int) {
			link := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea)
			if option == notSelectText {
				link.SetText("", false)
				return
			}
			if link.GetText() == "" {
				link.SetText(option, false)
				return
			}
			linkTagNames := strings.Split(link.GetText(), ",")
			linkTagNames = append(linkTagNames, option)
			linkTagNames = strutil.RemoveDuplicates(linkTagNames)
			link.SetText(strings.Join(linkTagNames, ","), false)
		}).
		AddButton("Save", func() {
			p.storeProject()
			tui.SetFocus("projectTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("projectTable")
		})
}

func (p *Project) setUpdateProjectForm(tui *service.TUI, projectID uint, projectName string, projectTagNames []string) {
	p.Form.Clear(true)
	p.ReadOnlyForm.Clear(true)

	p.ReadOnlyForm.AddTextArea("Selected Tags", "", 50, 5, 0, nil)
	tags := p.tagUC.GetAllNames()
	tags = append([]string{notSelectText}, tags...)
	p.Form.AddInputField("Project Name : ", projectName, 50, nil, nil).
		AddDropDown("Tags : ", tags, 0, func(option string, optionIndex int) {
			link := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea)
			if option == notSelectText {
				link.SetText("", false)
				return
			}
			if link.GetText() == "" {
				link.SetText(option, false)
				return
			}
			linkTagNames := strings.Split(link.GetText(), ",")
			linkTagNames = append(linkTagNames, option)
			linkTagNames = strutil.RemoveDuplicates(linkTagNames)
			link.SetText(strings.Join(linkTagNames, ","), false)
		}).
		AddButton("Update", func() {
			p.updateProject(projectID)
			tui.SetFocus("projectTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("projectTable")
		})
	p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea).SetText(strings.Join(projectTagNames, ","), false)
}

func (p *Project) formCapture(tui *service.TUI) {
	p.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlB:
			tui.SetFocus("projectTable")
		}
		return event
	})
}

func (p *Project) storeProject() error {
	projectName := p.Form.GetFormItemByLabel("Project Name : ").(*tview.InputField).GetText()
	projectTags := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea).GetText()
	if projectName == "" {
		return nil
	}

	var tagIDs []uint
	if projectTags != "" {
		tagNames := strings.Split(projectTags, ",")
		tags, err := p.tagUC.FindByNames(tagNames)
		if err == nil {
			for _, tag := range tags {
				tagIDs = append(tagIDs, tag.ID)
			}
		}
	}

	_, err := p.projectTypeUC.Create(projectName, tagIDs)
	if err != nil {
		return err
	}
	p.RestoreTable()

	return nil
}

func (p *Project) updateProject(projectID uint) {
	projectName := p.Form.GetFormItemByLabel("Project Name : ").(*tview.InputField).GetText()
	projectTags := p.ReadOnlyForm.GetFormItemByLabel("Selected Tags").(*tview.TextArea).GetText()
	if projectName == "" {
		return
	}

	var tagIDs []uint
	if projectTags != "" {
		tagNames := strings.Split(projectTags, ",")
		tags, err := p.tagUC.FindByNames(tagNames)
		if err == nil {
			for _, tag := range tags {
				tagIDs = append(tagIDs, tag.ID)
			}
		}
	}

	if err := p.projectTypeUC.Update(projectID, projectName, tagIDs); err != nil {
		p.errorHandler.ShowErrorWithErr(err, "projectTable")
		return
	}
	p.RestoreTable()
}

func (p *Project) setTable() {
	p.setTableHeader()
	p.setTableBody()
}

func (p *Project) setTableHeader() {
	for i, header := range projectHeader {
		tableCell := tview.NewTableCell(header).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorPurple).
			SetSelectable(false).
			SetExpansion(1)
		p.Table.SetCell(0, i, tableCell)
	}
}

func (p *Project) setTableBody() {
	projects, err := p.projectTypeUC.FindAllWithTags()
	if err != nil {
		p.errorHandler.ShowErrorWithErr(err, "projectTable")
		return
	}
	for i, project := range projects {
		p.Table.SetCell(i+1, 0,
			tview.NewTableCell(fmt.Sprint(project.ID)).
				SetAlign(tview.AlignCenter),
		)
		p.Table.SetCell(i+1, 1,
			tview.NewTableCell(project.Name).
				SetAlign(tview.AlignCenter),
		)
		tags := strings.Join(project.GetTagNames(), ",")
		p.Table.SetCell(i+1, 2,
			tview.NewTableCell(fmt.Sprint(tags)).
				SetAlign(tview.AlignCenter),
		)
	}
}

func (p *Project) tableCapture(tui *service.TUI) {
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'a':
				// store project
				p.setStoreProjectForm(tui)
				tui.SetFocus("projectForm")
			case 'u':
				// update project
				row, _ := p.Table.GetSelection()
				cell := p.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				id := cell.Text
				if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
					unitId := uint(intId)
					project, err := p.projectTypeUC.FindByID(unitId)
					if err != nil {
						p.errorHandler.ShowErrorWithErr(err, "projectTable")
						break
					}
					p.setUpdateProjectForm(tui, project.ID, project.Name, project.GetTagNames())
					tui.SetFocus("projectForm")
				}
			case 'd':
				// delete project
				row, _ := p.Table.GetSelection()
				cell := p.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				id := cell.Text

				var intId uint64
				isExist := false

				intId, _ = strconv.ParseUint(id, 10, 0)
				uintId := uint(intId)
				project, _ := p.projectTypeUC.FindByID(uintId)
				chronoWorks, err := p.chronoWorkUC.FindByProjectTypeID(project.ID)
				if err != nil {
					break
				}
				if len(chronoWorks) > 0 {
					isExist = true
				}
				var modal *tview.Modal
				if isExist {
					modal = tview.NewModal().
						SetText("Can't delete this project. Exist work that use this project.").
						AddButtons([]string{"Close"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							tui.DeleteModal()
							tui.SetFocus("projectTable")
							p.Table.ScrollToBeginning().Select(row, 0)
						})
				} else {
					modal = tview.NewModal().
						SetText("Are you sure you want to delete this project?").
						AddButtons([]string{"Yes", "No"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							if buttonLabel == "Yes" {
								if err := p.projectTypeUC.Delete(project.ID); err != nil {
									p.errorHandler.ShowErrorWithErr(err, "projectTable")
								}
								p.RestoreTable()
							}
							tui.DeleteModal()
							tui.SetFocus("projectTable")
							p.Table.ScrollToBeginning().Select(1, 0)
						})
				}
				tui.SetModal(modal)
				tui.SetFocus("modal")
			}
		}
		return event
	})
}
