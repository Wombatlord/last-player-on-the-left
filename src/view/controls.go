package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/wombatlord/last-player-on-the-left/src/audiopanel"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"unicode"
)

// Controller is the interface for views that are expressed
// In terms of tview.Primitive implementations. Implementations
// of this interface must pass their InputHandler method
// into tview.Primitive.SetInputCapture on initialisation or
// any controls they define will be ignored
type Controller interface {
	InputHandler(event *tcell.EventKey) *tcell.EventKey
}

// ReceiverController is an interface that gives the controller
// the capability to be notified if the global state changes
type ReceiverController interface {
	Controller
	domain.Receiver
}

// PanelStateAwareController is an interface that extends
// a Controller by requiring an implementation of OnUpdate
// which is periodically called with the audio panel state
type PanelStateAwareController interface {
	Controller
	audiopanel.PlayerStateSubscriber
}

// PanelStateAwareReceiverController is an interface that combines
// the functionality of both PanelStateAwareController and ReceiverController
type PanelStateAwareReceiverController interface {
	Controller
	domain.Receiver
	audiopanel.PlayerStateSubscriber
}

// Control is the abstraction of a keymapping
type Control func(event *tcell.EventKey) bool

var PlayPause Control = func(event *tcell.EventKey) bool {
	return event.Key() == tcell.KeyPause || unicode.ToLower(event.Rune()) == 'p'
}

var CycleFocus Control = func(event *tcell.EventKey) bool {
	return event.Key() == tcell.KeyTab
}

var FocusRight = func(event *tcell.EventKey) bool {
	return event.Key() == tcell.KeyRight
}

var FocusLeft = func(event *tcell.EventKey) bool {
	return event.Key() == tcell.KeyLeft
}

var SelectItem = func(event *tcell.EventKey) bool {
	return event.Key() == tcell.KeyEnter
}
