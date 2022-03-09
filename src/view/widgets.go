package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
)

var widgetLogger chan string

type CycleDirection int

const (
	CycleForward CycleDirection = iota
	CycleReverse
)

var nextIndex int

type Boxes struct {
	Left  *tview.List
	Right *tview.List
}

func InitState(conf app.Config, f *clients.RSSFeed) *Boxes {
	var err error
	if len(conf.Subs) == 0 {
		fmt.Println("No subs silly")
	} else {
		for _, url := range conf.Subs {
			f, err = clients.GetContent(url)
			if err != nil {
				panic(err)
			}
		}
	}
	return &Boxes{}
}

func (b *Boxes) FeedUrl() string {
	if b.Left == nil {
		return ""
	}
	_, url := b.Left.GetItemText(b.Left.GetCurrentItem())
	return url
}

func (b *Boxes) EpUrl() string {
	_, url := b.Right.GetItemText(b.Right.GetCurrentItem())
	return url
}

func (b *Boxes) Ordered() []*tview.List {
	return []*tview.List{b.Left, b.Right}
}

func (b *Boxes) Refresh() {
	PopulateMenu(LeftMenuProvider, b.Left.Clear())
	PopulateMenu(RightMenuProvider, b.Right.Clear())
}

type MenuProvider func() (string, []MenuItem)

var LeftMenuProvider MenuProvider
var RightMenuProvider MenuProvider

var State *Boxes

func AttachUI(guiApp *tview.Application) func() {
	widgetLogger = getLogger()
	widgetLogger <- "Attaching terminal GUI to App"
	State = &Boxes{}
	debug := tview.NewTextView()
	debug.SetTitle("DEBUG")

	State.Left = BuildMenu(LeftMenuProvider)
	State.Right = BuildMenu(RightMenuProvider)
	flex := insertPanes(State, debug)

	publishFunc := func(e *tcell.EventKey) *tcell.EventKey {
		e = overrides(e, State, guiApp, debug)
		return e
	}

	guiApp.
		SetInputCapture(publishFunc).
		SetRoot(flex, true).
		EnableMouse(true)

	quitFn := func() {
		guiApp.Stop()
	}

	return quitFn
}

func getLogger() chan string {
	return app.GetLogChan("view.logger")
}

func insertPanes(boxes *Boxes, debug *tview.TextView) *tview.Flex {
	flex := tview.NewFlex()

	// Set up the left pane content based on the presence of a debug panel
	var leftCol tview.Primitive = boxes.Left
	if debug != nil {
		leftCol = tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(boxes.Left, -1, 5, true).
			AddItem(debug, -1, 1, false)
	} else {
		leftCol = boxes.Left
	}

	flex.AddItem(leftCol, -1, 1, true)
	flex.AddItem(boxes.Right, -2, 1, false)

	return flex
}

type MenuItem struct {
	Name, Desc string
}

func BuildMenu(provider MenuProvider) *tview.List {
	menu := tview.NewList()
	PopulateMenu(provider, menu)

	return menu
}

func PopulateMenu(provider MenuProvider, menu *tview.List) {
	widgetLogger = getLogger()
	title, items := provider()
	for _, item := range items {
		widgetLogger <- fmt.Sprintf("menu item added %+v", item)
		menu.AddItem(
			item.Name, item.Desc, 0, nil,
		)
	}
	menu.SetTitle(title).SetBorder(true)
	menu.ShowSecondaryText(false)
}

func overrides(e *tcell.EventKey, boxes *Boxes, guiApp *tview.Application, debug *tview.TextView) *tcell.EventKey {
	widgetLogger = getLogger()
	switch e.Key() {
	case tcell.KeyEnter:
		mainText, _ := boxes.Left.GetItemText(boxes.Left.GetCurrentItem())
		widgetLogger <- mainText
	case tcell.KeyLeft:
		if guiApp.GetFocus() != boxes.Left {
			guiApp.SetFocus(boxes.Left)
		}
	case tcell.KeyRight:
		if guiApp.GetFocus() != boxes.Right {
			guiApp.SetFocus(boxes.Right)
		}
	}
	if debug != nil {
		debug.Clear()
		conf, _ := app.LoadConfig("config.yaml")
		debugOutput := []byte(
			fmt.Sprintf(
				"Last Event |=> %v\n"+
					"List Indices |=> %+v\n"+
					"Config |=> %+v\n"+
					"Next Focus Index |=> %d",
				e,
				struct{ Left, Right int }{
					boxes.Left.GetCurrentItem(),
					boxes.Right.GetCurrentItem(),
				},
				conf.Config,
				nextIndex,
			),
		)
		if _, err := debug.Write(debugOutput); err != nil {
			widgetLogger <- fmt.Sprintf("Error when writing to debug panel %+v", err)
		}
	}

	return e
}
