package view

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
)

type RootController struct {
	FlexController
	logger    chan string
	view      *tview.Flex
	gui       *tview.Application
	focusRing []tview.Primitive
}

func NewRootController(gui *tview.Application) *RootController {
	return &RootController{
		gui:    gui,
		logger: app.GetLogChan("RootController"),
	}
}

func (r *RootController) SetFocusRing(focusRing []tview.Primitive) {
	r.focusRing = focusRing
}

func (r *RootController) Attach(root *tview.Flex) {
	r.view = root
	root.SetInputCapture(r.InputHandler)
}

func (r *RootController) View() *tview.Flex {
	return r.view
}

func (r *RootController) Receive(_ app.State) {}

func (r *RootController) InputHandler(event *tcell.EventKey) *tcell.EventKey {
	var focusIndex int

	if event.Key() == tcell.KeyTab {
		for i := 0; i < 2; i++ {
			if r.focusRing[i] == r.gui.GetFocus() {
				focusIndex = i
				r.logger <- fmt.Sprintf("%+v", focusIndex)
			}
		}

		focusIndex += 1
		focusIndex = focusIndex % len(r.focusRing)
		r.gui.SetFocus(r.focusRing[focusIndex])
		return nil
	}

	if event.Key() == tcell.KeyRight {
		r.gui.SetFocus(r.focusRing[1])
		return nil
	}

	if event.Key() == tcell.KeyLeft {
		r.gui.SetFocus(r.focusRing[0])
		return nil
	}

	if event.Key() == tcell.KeyPause || unicode.ToLower(event.Rune()) == 'p' {
		lastplayer.FetchAudioPanel().PlayPause()
		return nil
	}

	return event

}
