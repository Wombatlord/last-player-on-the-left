package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/audiopanel"
	"log"
)

// RootController handles global controls
type RootController struct {
	Controller
	logger     *log.Logger
	lastPlayer *LastPlayer
}

// NewRootController initialises the RootController
func NewRootController(lastPlayer *LastPlayer) *RootController {
	r := &RootController{
		lastPlayer: lastPlayer,
		logger:     lastPlayer.GetLogger("RootController"),
	}
	r.lastPlayer.Views.Root.SetInputCapture(r.InputHandler)
	return r
}

// InputHandler implements the global controls. In some cases the events need to propagate
// through the hierarchy
func (r *RootController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	if CycleFocus(event) {
		focusIndex := (r.focusRingIndex() + 1) % len(r.lastPlayer.FocusRing)
		r.lastPlayer.SetFocus(r.lastPlayer.FocusRing[focusIndex])
		return nil
	}

	if FocusRight(event) {
		r.lastPlayer.SetFocus(r.lastPlayer.FocusRing[1])
		return nil
	}

	if FocusLeft(event) {
		r.lastPlayer.SetFocus(r.lastPlayer.FocusRing[0])
		return nil
	}

	if PlayPause(event) {
		audiopanel.FetchAudioPanel().PlayPause()
		return event
	}

	return event

}

// focusRingIndex returns the index of the currently focussed view relative to the
// LastPlayer.FocusRing
func (r *RootController) focusRingIndex() int {
	var focusIndex int
	for i := 0; i < 2; i++ {
		if r.lastPlayer.FocusRing[i] == r.lastPlayer.GetFocus() {
			focusIndex = i
			r.logger.Printf("%+v", focusIndex)
		}
	}
	return focusIndex
}

// focusedView returns the currently focussed view
func (r *RootController) focusedView() tview.Primitive {
	return r.lastPlayer.FocusRing[r.focusRingIndex()]
}
