package view

import (
	"github.com/rivo/tview"
)

type Controllers struct {
	FeedMenu      *FeedsMenuController
	EpisodeMenu   *EpisodeMenuController
	RootContoller *RootContoller
}

// Build returns the tview app, implement any additions to
// the user interface adding new primitives to the app hierarchy
// in this file
func Build(gui *tview.Application, controllers Controllers) *tview.Application {
	mainFlex := MainFlex(controllers)
	feedColumn := FeedColumn()
	episodeMenu := EpisodeMenu(controllers)
	feedMenu := FeedMenu(controllers)
	debugPanel := DebugPanel()

	feedColumn.AddItem(feedMenu, -1, 5, true)
	feedColumn.AddItem(debugPanel, -1, 1, false)

	mainFlex.AddItem(feedColumn, -1, 1, true)
	mainFlex.AddItem(episodeMenu, -1, 1, true)

	gui.SetRoot(mainFlex, true)

	focusRing := []tview.Primitive{feedMenu, episodeMenu}
	controllers.RootContoller.SetFocusRing(focusRing)


	return gui
}

func MainFlex(controllers Controllers) *tview.Flex {
	mainFlex := tview.NewFlex()
	controllers.RootContoller.Attach(mainFlex)

	return mainFlex
}

func FeedColumn() *tview.Flex {
	feedColumnFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	return feedColumnFlex
}

func EpisodeMenu(controllers Controllers) *tview.List {
	episodeMenuView := tview.NewList()
	controllers.EpisodeMenu.Attach(episodeMenuView)

	episodeMenuView.SetBorder(true)
	return episodeMenuView
}

func FeedMenu(controllers Controllers) *tview.List {
	feedMenu := tview.NewList()
	controllers.FeedMenu.Attach(feedMenu)

	feedMenu.SetBorder(true)
	return feedMenu
}

func DebugPanel() *tview.TextView {
	view := tview.NewTextView()
	return view
}
