package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"log"
)

// FeedsMenuController manages the feeds menu, it synchronises the
// current feed with the feed selection in the ui
type FeedsMenuController struct {
	Controller
	lastPlayer *LastPlayer
	logger     *log.Logger
}

// NewFeedsController initialises the FeedsMenuController
func NewFeedsController(application *LastPlayer) *FeedsMenuController {
	f := &FeedsMenuController{
		logger:     application.GetLogger("FeedsMenuController"),
		lastPlayer: application,
	}
	application.Views.FeedMenu.SetInputCapture(f.InputHandler)
	for _, sub := range app.LoadedConfig.Subs {
		application.Views.FeedMenu.AddItem(sub.Alias, sub.Url, 0, nil)
	}

	return f
}

// selectFeed updates the UI with the new feed as selected by the user
func (f *FeedsMenuController) selectFeed() {
	f.logger.Printf(
		"Pushing feed index %d to state",
		f.lastPlayer.Views.FeedMenu.GetCurrentItem(),
	)

	go f.lastPlayer.QueueUpdateDraw(func() {
		f.lastPlayer.State.Feed = f.getFeed()
		f.lastPlayer.State.FeedIndex = f.lastPlayer.Views.FeedMenu.GetCurrentItem()
	})
}

// getFeed returns the *clients.RSSFeed associated to the current selection in the menu
func (f *FeedsMenuController) getFeed() *clients.RSSFeed {
	index := f.lastPlayer.Views.FeedMenu.GetCurrentItem()
	_, url := f.lastPlayer.Views.FeedMenu.GetItemText(index)
	feed, err := clients.GetContent(url)
	if err != nil {
		log.Fatal(err)
	}

	return feed
}

// InputHandler invokes selectFeed on capturing a tcell.KeyEnter keypress
func (f *FeedsMenuController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if SelectItem(event) {
		f.selectFeed()
		f.lastPlayer.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
	}
	return event
}
