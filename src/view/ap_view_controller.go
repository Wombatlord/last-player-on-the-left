package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/wombatlord/last-player-on-the-left/src/audiopanel"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"log"
)

// APViewController manages the updating of the tview.TextView that shows the current
// playing episode
type APViewController struct {
	PanelStateAwareReceiverController
	lastPlayer     *LastPlayer
	logger         *log.Logger
	playingEpisode *clients.Item
}

// NewAPViewController initialises the APViewController and provides an interface
// for dependency injection
func NewAPViewController(lastPlayer *LastPlayer) *APViewController {
	a := &APViewController{logger: lastPlayer.GetLogger("APViewController"), lastPlayer: lastPlayer}
	return a
}

// OnUpdate implements the audiopanel.PlayerStateSubscriber interface
func (a *APViewController) OnUpdate() {
	state := a.lastPlayer.AudioPanel.GetPlayerState()
	a.RenderState(state)
}

// RenderState updates the TextView that the APViewController controls
func (a *APViewController) RenderState(state audiopanel.PlayerState) {
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
	a.lastPlayer.Views.APView.SetText(playerStatus)
}

// Receive updates the internal tracking of the currently playing episode, if it has
// changed then the view is re-rendered
func (a *APViewController) Receive(state domain.State) {
	a.logger.Printf("Received state %+v", state)
	if *a.playingEpisode != *state.PlayingEpisode {
		a.playingEpisode = state.PlayingEpisode
		a.RenderState(a.lastPlayer.AudioPanel.GetPlayerState())
	}
}

// InputHandler is used here to rerender the view with the updated player state on capture
// of the 'Play/Pause' control input
func (a *APViewController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if PlayPause(event) {
		a.RenderState(a.lastPlayer.AudioPanel.GetPlayerState())
		return nil
	}
	return event
}
