package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"log"
)

// FeedsMenuController manages the feeds menu, it synchronises the
// current feed with the feed selection in the ui
type FeedsMenuController struct {
	BaseMenuController
	view                  *tview.List
	hightlightedFeedIndex int
	feed                  *clients.RSSFeed
	logger                chan string
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
	f.hightlightedFeedIndex = list.GetCurrentItem()
	f.view = list
}

func (f *FeedsMenuController) OnSelectionChange(
	index int,
	mainText string,
	secondaryText string,
	shortcut rune,
) {
	var err error
	if secondaryText != "" {
		f.hightlightedFeedIndex = index
	}
	if err != nil {
		f.logger <- fmt.Sprintf("Error occurred while attempting to retrieve the feed: %s", err.Error())
	}
	//f.selectFeed(index)
}

// selectFeed updates the UI with the new feed as selected by the user
func (f *FeedsMenuController) selectFeed() {
	f.logger <- fmt.Sprintf("Pushing feed index %d to state", f.hightlightedFeedIndex)
	index := f.view.GetCurrentItem()
	var err error
	_, url := f.view.GetItemText(index)
	f.feed, err = clients.GetContent(url)
	if err != nil {
		log.Fatal(err)
	}

	manager := domain.NewManager()
	manager.QueueTransform(
		func(state domain.State) domain.State {
			state.FeedIndex = f.hightlightedFeedIndex
			state.Feed = f.feed
			return state
		},
	)
	manager.Commit()
}

func (f *FeedsMenuController) Receive(s domain.State) {
	f.logger <- fmt.Sprintf("Received state: %+v", s)
}

func (f *FeedsMenuController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		f.selectFeed()
	}
	return event
}
