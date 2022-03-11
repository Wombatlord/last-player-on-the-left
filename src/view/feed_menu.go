package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
)

// FeedsMenuController manages the feeds menu, it synchronises the
// current feed with the feed selection in the ui
type FeedsMenuController struct {
	BaseMenuController
	view      *tview.List
	feedIndex int
	logger    chan string
}

func NewFeedsController() *FeedsMenuController {
	return &FeedsMenuController{logger: app.GetLogChan("FeedsMenuController")}
}

func (f *FeedsMenuController) Attach(list *tview.List) {
	list.SetChangedFunc(f.OnSelectionChange)
	list.SetInputCapture(f.InputHandler)
	for _, sub := range app.LoadedConfig.Subs {
		list.AddItem(sub.Alias, sub.Url, 0, nil)
	}
	f.view = list
}

func (f *FeedsMenuController) OnSelectionChange(
	index int,
	mainText string,
	secondaryText string,
	shortcut rune,
) {
	f.feedIndex = index
	//f.updateFeedIndex(index)
}

func (f *FeedsMenuController) updateFeedIndex(index int) {
	manager := app.NewManager()

	manager.QueueTransform(
		func(state app.State) app.State {
			state.FeedIndex = index
			return state
		},
	)
	manager.Commit()
}

func (f *FeedsMenuController) Receive(s app.State) {
	f.logger <- fmt.Sprintf("Received state: %+v", s)
}

func (f *FeedsMenuController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		f.updateFeedIndex(f.view.GetCurrentItem())
	}
	return event
}
