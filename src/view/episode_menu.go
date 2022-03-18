package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"log"
)

// EpisodeMenuController Handles input captured from and updates to be
// displayed in the episode menu
type EpisodeMenuController struct {
	ReceiverController
	feed           *clients.RSSFeed
	feedIndex      int
	playingEpisode *clients.Item
	lastPlayer     *LastPlayer
	logger         *log.Logger
}

// NewEpisodeMenuController Initialises the EpisodeMenuController
func NewEpisodeMenuController(lastPlayer *LastPlayer) *EpisodeMenuController {
	e := &EpisodeMenuController{
		lastPlayer: lastPlayer,
		feedIndex:  domain.NoItem,
		logger:     lastPlayer.GetLogger("EpisodeMenuController"),
	}
	lastPlayer.Views.EpisodeMenu.SetInputCapture(e.InputHandler)
	return e
}

// playEpisode retrieves the appropriate feed item and uses
// the url in its Enclosure to pass to panel.PlayFromUrl which initiates
// audio playback. It then queues a function to make the information about the
// currently playing episode available to the rest of the application
// via the global state.
func (e *EpisodeMenuController) playEpisode() {
	episodeIndex := e.lastPlayer.Views.EpisodeMenu.GetCurrentItem()
	e.playingEpisode = &e.lastPlayer.State.Feed.Channel[0].Item[episodeIndex]

	panel := e.lastPlayer.AudioPanel
	panel.PlayFromUrl(e.playingEpisode.Enclosure.Url)
	go e.lastPlayer.QueueUpdateDraw(func() {
		e.lastPlayer.State.PlayingEpisode = e.playingEpisode
	})
}

// Receive is looking out for changes to the feed index
func (e *EpisodeMenuController) Receive(state domain.State) {
	e.logger.Printf("Receive called with state %v", state)
	e.update(state)
}

// update sets the view state so that it us redrawn on the next
// application draw cycle
func (e *EpisodeMenuController) update(state domain.State) {
	if e.feedIndex == domain.NoItem || e.feedIndex != state.FeedIndex {
		alias := "NoItem"
		if state.FeedIndex > domain.NoItem {
			alias = e.lastPlayer.Config.Subs[state.FeedIndex].Alias
		}
		e.logger.Printf("Feed changed to %s, redrawing menu", alias)
		e.feedIndex = state.FeedIndex
		e.lastPlayer.Views.EpisodeMenu.Clear()
		for _, item := range e.lastPlayer.State.Feed.Channel[0].Item {
			e.lastPlayer.Views.EpisodeMenu.AddItem(item.Title, item.Enclosure.Url, ' ', nil)
		}
	}
}

// InputHandler implements the user input side of the controller interface
func (e *EpisodeMenuController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if SelectItem(event) {
		e.playEpisode()
		return nil
	}

	return event
}
