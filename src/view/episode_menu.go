package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
)

type EpisodeMenuController struct {
	BaseMenuController
	feed      *clients.RSSFeed
	feedIndex int
	view      *tview.List
	logger    chan string
}

func NewEpisodeMenuController() *EpisodeMenuController {
	return &EpisodeMenuController{feedIndex: -1, logger: app.GetLogChan("EpisodeMenuController")}
}

func (e *EpisodeMenuController) Attach(list *tview.List) {
	list.SetChangedFunc(e.OnSelectionChange)
	list.SetInputCapture(e.InputHandler)
	e.view = list
}

func (e *EpisodeMenuController) OnSelectionChange(
	index int,
	_ string,
	_ string,
	_ rune,
) {
	manager := app.NewManager()
	manager.QueueTransform(
		func(state app.State) app.State {
			state.EpisodeIndex = index
			return state
		},
	)

	manager.Commit()
}

// Receive is looking out for changes to the feed index
func (e *EpisodeMenuController) Receive(s app.State) {
	e.logger <- fmt.Sprintf("Received state: %+v", s)
	if e.feedIndex == -1 || e.feedIndex != s.FeedIndex {
		e.feed, _ = clients.GetContent(app.LoadedConfig.Subs[s.FeedIndex].Url)
		e.feedIndex = s.FeedIndex
		e.view.Clear()
		for _, item := range e.feed.Channel[0].Item {
			e.view.AddItem(item.Title, item.Enclosure.Url, ' ', nil)
		}
	}
}

func (e *EpisodeMenuController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		episodeIndex := e.view.GetCurrentItem()
		lastplayer.FetchAudioPanel().PlayFromUrl(e.feed.Channel[0].Item[episodeIndex].Enclosure.Url)
		return nil
	}

	return event
}
