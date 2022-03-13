package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
	"time"
)

type APViewController struct {
	TextViewController
	lastplayer.PlayerStateSubscriber
	view   *tview.TextView
	logger chan string
	ticker *time.Ticker
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
	a.view.SetText(fmt.Sprintf("STATE: %+v", state))
}

func (a *APViewController) View() *tview.TextView {
	return a.view
}

func (a *APViewController) Receive(state app.State) {
	a.logger <- fmt.Sprintf("Received state %+v", state)
}

func (a *APViewController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	return event
}
