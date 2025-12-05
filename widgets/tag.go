package widgets

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/internal/usecase"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/rivo/tview"
)

var (
	tagHeader = []string{
		"ID",
		"Name",
	}
)

type Tag struct {
	Layout *tview.Grid
	Form   *tview.Form
	Table  *tview.Table
	tagUC  *usecase.TagUseCase
}

func NewTag(tagUC *usecase.TagUseCase) *Tag {
	return &Tag{
		Layout: tview.NewGrid().
			SetRows(10, 0).
			SetColumns(0).
			SetBorders(true),
		Form: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple).
			SetFieldTextColor(tcell.ColorGray).
			SetFieldBackgroundColor(tcell.ColorWhite),
		Table: tview.NewTable().
			SetSelectable(true, false).
			SetFixed(1, 1),
		tagUC: tagUC,
	}
}

func (t *Tag) GenerateInitTag(tui *service.TUI) *Tag {
	t.setStoreTagForm(tui)
	t.restoreTable()

	t.Layout.AddItem(t.Form, 0, 0, 1, 1, 0, 0, false)
	t.Layout.AddItem(t.Table, 1, 0, 1, 1, 0, 0, true)

	t.tableCapture(tui)
	t.formCapture(tui)
	return t
}

func (t *Tag) tableCapture(tui *service.TUI) {
	t.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'a':
				t.setStoreTagForm(tui)
				tui.SetFocus("tagForm")
			case 'u':
				row, _ := t.Table.GetSelection()
				cell := t.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				id := cell.Text
				if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
					unitId := uint(intId)
					tag, err := t.tagUC.FindByID(unitId)
					if err != nil {
						break
					}
					t.setUpdateTagForm(tui, tag.ID, tag.Name)
					tui.SetFocus("tagForm")
				}

			}
		}
		return event
	})
}

func (t *Tag) formCapture(tui *service.TUI) {
	t.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlB:
			tui.SetFocus("tagTable")
		}
		return event
	})
}

func (t *Tag) setStoreTagForm(tui *service.TUI) {
	t.Form.Clear(true)
	t.Form.
		AddInputField("Name", "", 50, nil, nil).
		AddButton("Create", func() {
			err := t.storeTag()
			if err != nil {
				return
			}
			tui.SetFocus("tagTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("tagTable")
		})
}

func (t *Tag) setUpdateTagForm(tui *service.TUI, tagID uint, tagName string) {
	t.Form.Clear(true)
	t.Form.
		AddInputField("Name", tagName, 50, nil, nil).
		AddButton("Update", func() {
			err := t.updateTag(tagID)
			if err != nil {
				return
			}
			tui.SetFocus("tagTable")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("tagTable")
		})
}

func (t *Tag) storeTag() error {
	tagName := t.Form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	if tagName == "" {
		return fmt.Errorf("tag name is empty")
	}
	_, err := t.tagUC.Create(tagName)
	if err != nil {
		return err
	}
	t.restoreTable()
	return nil
}

func (t *Tag) updateTag(tagID uint) error {
	tagName := t.Form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	if tagName == "" {
		return fmt.Errorf("tag name is empty")
	}
	if err := t.tagUC.Update(tagID, tagName); err != nil {
		return err
	}
	t.restoreTable()
	return nil
}

func (t *Tag) setTable() {
	t.setTableHeader()
	t.setTableBody()
}

func (t *Tag) setTableHeader() {
	for i, header := range tagHeader {
		t.Table.SetCell(0, i,
			tview.NewTableCell(header).
				SetAlign(tview.AlignLeft).
				SetTextColor(tcell.ColorWhite).
				SetBackgroundColor(tcell.ColorPurple).
				SetSelectable(false))
	}
}

func (t *Tag) setTableBody() {
	tags, err := t.tagUC.FindAll()
	if err != nil {
		return
	}
	for i, tag := range tags {
		t.Table.SetCell(i+1, 0,
			tview.NewTableCell(fmt.Sprint(tag.ID)).
				SetAlign(tview.AlignCenter).
				SetExpansion(0))
		t.Table.SetCell(i+1, 1,
			tview.NewTableCell(tag.Name).
				SetAlign(tview.AlignLeft).
				SetExpansion(1))
	}
}

func (t *Tag) restoreTable() {
	t.Table.Clear()
	t.setTable()
}
