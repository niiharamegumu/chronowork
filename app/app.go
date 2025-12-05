package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/container"
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

	// Initialize DI container
	c := container.New(db.DB)

	// Get relative days from setting
	setting, err := c.SettingUC.Get()
	if err != nil {
		return err
	}
	relativeDays := int(setting.RelativeDate)

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("ChronoWork")

	mainTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("Today is %s (%v)", time.Now().Format("2006/01/02"), time.Now().Weekday())).SetTextColor(tcell.ColorPurple)

	// Initialize widgets with use cases
	timer := widgets.NewTimer(c.ChronoWorkUC)
	err = timer.CheckActiveTracking(tui)
	if err != nil {
		return err
	}

	work := widgets.NewWork(c.ChronoWorkUC, c.SettingUC)
	work, err = work.GenerateInitWork(tui, relativeDays)
	if err != nil {
		return err
	}

	form := widgets.NewForm(c.ChronoWorkUC, c.ProjectTypeUC)
	form = form.GenerateInitForm(tui, work, relativeDays)

	// add page
	// setting page
	settingWidget := widgets.NewSetting(c.SettingUC)
	settingWidget.GenerateInitSetting(tui)
	tui.SetMainPage("setting", settingWidget.Form, false)
	if err = tui.SetWidget("settingForm", settingWidget.Form); err != nil {
		return err
	}

	// project page
	project := widgets.NewProject(c.ProjectTypeUC, c.TagUC, c.ChronoWorkUC)
	tui.SetMainPage("project", project.Layout, false)
	if err = tui.SetWidget("projectForm", project.Form); err != nil {
		return err
	}
	if err = tui.SetWidget("projectTable", project.Table); err != nil {
		return err
	}
	project.GenerateInitProject(tui)

	// tag page
	tagPage := widgets.NewTag(c.TagUC)
	tagPage.GenerateInitTag(tui)
	tui.SetMainPage("tag", tagPage.Layout, false)
	if err = tui.SetWidget("tagForm", tagPage.Form); err != nil {
		return err
	}
	if err = tui.SetWidget("tagTable", tagPage.Table); err != nil {
		return err
	}

	// export page
	export := widgets.NewExport(c.ChronoWorkUC, c.SettingUC)
	export.GenerateInitExport(tui)
	tui.SetMainPage("export", export.Form, false)
	if err = tui.SetWidget("exportForm", export.Form); err != nil {
		return err
	}

	menu := widgets.NewMenu(c.SettingUC)
	menu = menu.GenerateInitMenu(tui, work, settingWidget, project)

	tui.SetHeader(header, false)
	tui.SetMenu(menu.List, false)
	tui.SetWork(mainTitle, form.Form, timer.Wrapper, work.Table, true) // default focus
	work.TableCapture(tui, form, timer, relativeDays)
	form.FormCapture(tui)

	tui.GlobalKeyActions()
	if err = tui.App.SetRoot(tui.Grid, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
