package widgets

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/internal/usecase"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/niiharamegumu/chronowork/util/timeutil"
	"github.com/rivo/tview"
)

type Export struct {
	Form         *tview.Form
	chronoWorkUC *usecase.ChronoWorkUseCase
	settingUC    *usecase.SettingUseCase
}

func NewExport(chronoWorkUC *usecase.ChronoWorkUseCase, settingUC *usecase.SettingUseCase) *Export {
	return &Export{
		Form: tview.NewForm().
			SetLabelColor(tcell.ColorPurple),
		chronoWorkUC: chronoWorkUC,
		settingUC:    settingUC,
	}
}

func (e *Export) GenerateInitExport(tui *service.TUI) {
	e.Form.AddButton("Export", func() {
		e.export()
		e.ReStore(tui)
		tui.SetFocus("menu")
	}).
		AddButton("Cancel", func() {
			e.ReStore(tui)
			tui.SetFocus("menu")
		})
}

func (e *Export) ReStore(tui *service.TUI) {
	e.Form.Clear(true)
	e.GenerateInitExport(tui)
}

func (e *Export) export() {
	setting, err := e.settingUC.Get()
	if err != nil {
		log.Println(err)
		return
	}
	path := setting.DownloadPath

	chronoWorks, err := e.chronoWorkUC.GetAll("id", 0)
	if err != nil {
		log.Println(err)
		return
	}
	if len(chronoWorks) < 1 {
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println(err)
		return
	}
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	baseFileName := fmt.Sprintf("%schrono_works", path)
	timestamp := time.Now().Format("20060102150405")
	exportPath := fmt.Sprintf("%s_%s.csv", baseFileName, timestamp)

	f, err := os.Create(exportPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{"ID", "Title", "ProjectName", "TagName", "Date", "Time"}
	if err := w.Write(header); err != nil {
		log.Println(err)
		return
	}

	for _, c := range chronoWorks {
		projectName := ""
		tagName := ""
		if c.ProjectType != nil {
			projectName = c.ProjectType.Name
		}
		if c.Tag != nil {
			tagName = c.Tag.Name
		}
		record := []string{
			strconv.Itoa(int(c.ID)),
			c.Title,
			projectName,
			tagName,
			c.CreatedAt.Format("2006/01/02"),
			timeutil.FormatTime(c.TotalSeconds),
		}
		if err := w.Write(record); err != nil {
			log.Println(err)
			continue
		}
	}
}
