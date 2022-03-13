package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
)

type EpisodeMenuController struct {
	BaseMenuController
	feed           *clients.RSSFeed
	feedIndex      int
	playingEpisode *clients.Item
	view           *tview.List
	logger         chan string
}

func NewEpisodeMenuController() *EpisodeMenuController {
	return &EpisodeMenuController{feedIndex: domain.NoItem, logger: app.GetLogChan("EpisodeMenuController")}
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
	e.highlightEpisode()
}

func (e *EpisodeMenuController) highlightEpisode() {
	manager := domain.NewManager()
	manager.QueueTransform(
		func(state domain.State) domain.State {
			state.EpisodeIndex = e.view.GetCurrentItem()
			return state
		},
	)

	manager.Commit()
}

func (e *EpisodeMenuController) playEpisode() {
	episodeIndex := e.view.GetCurrentItem()
	e.playingEpisode = &e.feed.Channel[0].Item[episodeIndex]

	panel := lastplayer.FetchAudioPanel()
	panel.PlayFromUrl(e.playingEpisode.Enclosure.Url)

	manager := domain.NewManager()
	manager.QueueTransform(
		func(state domain.State) domain.State {
			state.PlayingEpisode = e.playingEpisode
			return state
		},
	)
	manager.Commit()
}

// Receive is looking out for changes to the feed index
func (e *EpisodeMenuController) Receive(s domain.State) {
	e.feed = s.Feed
	e.logger <- fmt.Sprintf("Received state: %+v", s)
	if e.feedIndex == domain.NoItem || e.feedIndex != s.FeedIndex {
		e.feedIndex = s.FeedIndex
		e.view.Clear()
		for _, item := range e.feed.Channel[0].Item {
			e.view.AddItem(item.Title, item.Enclosure.Url, ' ', nil)
		}
	}
}

func (e *EpisodeMenuController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		e.playEpisode()
		return nil
	}

	return event
}
