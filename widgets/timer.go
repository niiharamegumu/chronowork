package widgets

import (
	"context"
	"time"
	
	"github.com/niiharamegumu/chronowork/db"
	"github.com/niiharamegumu/chronowork/models"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/niiharamegumu/chronowork/util/timeutil"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Timer struct {
	Wrapper     *tview.Grid
	Time        *tview.TextView
	Title       *tview.TextView
	CreatedDate *tview.TextView
	ProjectName *tview.TextView
	TagName     *tview.TextView
	StartTime   time.Time
	cancelCtx   context.Context
	cancelFunc  context.CancelFunc
}

func NewTimer() *Timer {
	time := tview.NewTextView().
		SetLabel("Timer : ").
		SetTextColor(tcell.ColorPurple).
		SetText("00:00:00")
	title := tview.NewTextView().
		SetTextColor(tcell.ColorPurple).
		SetLabel("TItle : ")
	CreatedDate := tview.NewTextView().
		SetTextColor(tcell.ColorPurple).
		SetLabel("Created Date : ")
	projectName := tview.NewTextView().
		SetTextColor(tcell.ColorPurple).
		SetLabel("Project Name : ")
	tagName := tview.NewTextView().
		SetTextColor(tcell.ColorPurple).
		SetLabel("Tag Name : ")
	timer := &Timer{
		Wrapper: tview.NewGrid().
			SetRows(0, 1, 1, 1, 0).
			SetColumns(0).
			AddItem(time, 0, 0, 1, 1, 0, 0, false).
			AddItem(title, 1, 0, 1, 1, 0, 0, false).
			AddItem(CreatedDate, 2, 0, 1, 1, 0, 0, false).
			AddItem(projectName, 3, 0, 1, 1, 0, 0, false).
			AddItem(tagName, 4, 0, 1, 1, 0, 0, false),
		Time:        time,
		Title:       title,
		CreatedDate: CreatedDate,
		ProjectName: projectName,
		TagName:     tagName,
	}
	return timer
}

func (t *Timer) CheckActiveTracking(tui *service.TUI) error {
	trackingChronoWork, err := models.FindTrackingChronoWorks(db.DB)
	if err != nil {
		return err
	}
	if len(trackingChronoWork) > 0 {
		t.SetStartTimer(trackingChronoWork[0].StartTime)
		t.SetCalculateSeconds(tui)
		t.SetTimerText(trackingChronoWork[0])
	}
	return nil
}

func (t *Timer) SetTimerText(c models.ChronoWork) {
	t.Title.SetText(c.Title)
	t.CreatedDate.SetText(c.CreatedAt.Format("2006-01-02 "))
	t.ProjectName.SetText(c.ProjectType.Name)
	t.TagName.SetText(c.Tag.Name)
}

func (t *Timer) SetStartTimer(startTime time.Time) {
	t.StartTime = startTime
}

func (t *Timer) SetCalculateSeconds(tui *service.TUI) {
	t.cancelCtx, t.cancelFunc = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-t.cancelCtx.Done():
				return
			default:
				seconds := int(time.Since(t.StartTime).Seconds())
				tui.App.QueueUpdateDraw(func() {
					t.Time.SetText(timeutil.FormatTime(seconds))
				})
				time.Sleep(time.Second)
			}
		}
	}()
}

func (t *Timer) StopCalculateSeconds() {
	t.cancelFunc()
}

func (t *Timer) ResetSetText() {
	t.Time.SetText("00:00:00")
	t.Title.SetText("")
	t.CreatedDate.SetText("")
	t.ProjectName.SetText("")
	t.TagName.SetText("")
}
