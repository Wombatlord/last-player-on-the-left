package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
)

// Controller is the interface for views that are expressed
// In terms of tview.Primitive implementations
type Controller interface {
	app.Receiver
	InputHandler(event *tcell.EventKey) *tcell.EventKey
}

// MenuController must be implemented for any controller
// that will attach to a tview.List
type MenuController interface {
	Controller
	Attach(view *tview.List)
	View() *tview.List
	OnSelectionChange(
		index int,
		mainText string,
		secondaryText string,
		shortcut rune,
	)
}

// BaseMenuController sets up standard implementations of Attach and View.
// Override them in your controller implementation to set up custom behavior
// on set/get
type BaseMenuController struct {
	MenuController
	view *tview.List
}

func (f *BaseMenuController) Attach(list *tview.List) {
	list.SetChangedFunc(f.OnSelectionChange)
	list.SetInputCapture(f.InputHandler)
	f.view = list
}

func (f *BaseMenuController) View() *tview.List {
	return f.view
}

// FlexController is a Controller that can be attached to a flex
type FlexController interface {
	Controller
	Attach(view *tview.Flex)
	View() *tview.Flex
}
