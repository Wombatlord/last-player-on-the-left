package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
	"time"
)

type APViewController struct {
	TextViewController
	lastplayer.PlayerStateSubscriber
	view           *tview.TextView
	logger         chan string
	ticker         *time.Ticker
	feed           *clients.RSSFeed
	playingEpisode *clients.Item
}

func NewAPViewController() *APViewController {
	return &APViewController{logger: app.GetLogChan("APViewController")}
}

func (a *APViewController) Attach(textView *tview.TextView) {
	a.view = textView
}

func (a *APViewController) OnUpdate() {
	state := lastplayer.FetchAudioPanel().GetPlayerState()
	a.logger <- fmt.Sprintf("OnUpdate called with state %+v", state)
	title := ""
	description := ""
	if a.playingEpisode != nil {
		title = a.playingEpisode.Title
		description = a.playingEpisode.Description
	}
	playStatus := map[bool]string{true: string(''), false: string('')}
	playerStatus := fmt.Sprintf(
		"%s\n%s %s/%s\n%s",
		title,
		playStatus[state.Playing],
		state.Position,
		state.Length,
		description,
	)
	a.view.SetText(playerStatus)
}

func (a *APViewController) View() *tview.TextView {
	return a.view
}

func (a *APViewController) Receive(state domain.State) {
	a.logger <- fmt.Sprintf("Received state %+v", state)
	a.feed = state.Feed
	a.playingEpisode = state.PlayingEpisode
}

func (a *APViewController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	return event
}
