package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/db"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/niiharamegumu/chronowork/widgets"
	"github.com/rivo/tview"
)

var (
	tui = service.NewTUI()
)

func init() {
	err := db.ConnectDB()
	if err != nil {
		fmt.Println("database connection error", err)
		os.Exit(1)
	}
	log.Println("database connection success")
}

func Execute() {
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Println("error closing database", err)
		}
		log.Println("database connection closed")
	}()

	// ==============================
	// [SEEDER] create test data
	// ==============================
	// if err := db.CreateTestData(db.DB); err != nil {
	// 	log.Println("error creating test data", err)
	// }

	if err := initialSetting(); err != nil {
		log.Println("error", err)
		os.Exit(1)
	}
}

func initialSetting() error {
	var err error

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("ChronoWork")

	mainTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("Today is %s (%v)", time.Now().Format("2006/01/02"), time.Now().Weekday())).SetTextColor(tcell.ColorPurple)
	timer := widgets.NewTimer()
	err = timer.CheckActiveTracking(tui)
	if err != nil {
		return err
	}
	work := widgets.NewWork()
	work, err = work.GenerateInitWork(tui)
	if err != nil {
		return err
	}
	form := widgets.NewForm()
	form = form.GenerateInitForm(tui, work)

	// add page
	// setting page
	setting := widgets.NewSetting()
	setting.GenerateInitSetting(tui)
	tui.SetMainPage("setting", setting.Form, false)
	if err = tui.SetWidget("settingForm", setting.Form); err != nil {
		return err
	}
	// project page
	project := widgets.NewProject()
	tui.SetMainPage("project", project.Layout, false)
	if err = tui.SetWidget("projectForm", project.Form); err != nil {
		return err
	}
	if err = tui.SetWidget("projectTable", project.Table); err != nil {
		return err
	}
	project.GenerateInitProject(tui)
	// tag page
	tagPage := widgets.NewTag()
	tagPage.GenerateInitTag(tui)
	tui.SetMainPage("tag", tagPage.Layout, false)
	if err = tui.SetWidget("tagForm", tagPage.Form); err != nil {
		return err
	}
	if err = tui.SetWidget("tagTable", tagPage.Table); err != nil {
		return err
	}
	// export page
	export := widgets.NewExport()
	export.GenerateInitExport(tui)
	tui.SetMainPage("export", export.Form, false)
	if err = tui.SetWidget("exportForm", export.Form); err != nil {
		return err
	}

	menu := widgets.NewMenu()
	menu = menu.GenerateInitMenu(tui, work, setting, project)

	tui.SetHeader(header, false)
	tui.SetMenu(menu.List, false)
	tui.SetWork(mainTitle, form.Form, timer.Wrapper, work.Table, true) // default focus
	work.TableCapture(tui, form, timer)
	form.FormCapture(tui)

	tui.GlobalKeyActions()
	if err = tui.App.SetRoot(tui.Grid, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
