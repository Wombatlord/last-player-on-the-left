package view

import (
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
}

func NewFeedsController() *FeedsMenuController {
	return &FeedsMenuController{}
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
	manager := app.NewManager()
	defer manager.Commit()

	manager.QueueTransform(
		func(state app.State) app.State {
			state.FeedIndex = index
			return state
		},
	)
}

func (f *FeedsMenuController) Receive(s app.State) {}
func (f *FeedsMenuController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	return event
}
