package view

import (
	"github.com/rivo/tview"
)

type Controllers struct {
	FeedMenu    *FeedsMenuController
	EpisodeMenu *EpisodeMenuController
}

func Build(controllers Controllers) *tview.Application {
	gui := tview.NewApplication()
	gui.SetRoot(MainFlex(controllers), true)
	return gui
}

func MainFlex(controllers Controllers) *tview.Flex {
	mainFlex := tview.NewFlex()
	mainFlex.AddItem(FeedColumn(controllers), -1, 1, true)
	mainFlex.AddItem(EpisodeMenu(controllers), -1, 1, true)
	return mainFlex
}

func FeedColumn(controllers Controllers) *tview.Flex {
	feedColumnFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	feedColumnFlex.AddItem(FeedMenu(controllers), -1, 5, true)
	feedColumnFlex.AddItem(DebugPanel(controllers), -1, 1, false)
	return feedColumnFlex
}

func EpisodeMenu(controllers Controllers) *tview.List {
	episodeMenuView := tview.NewList()
	episodeMenuView.SetBorder(true)
	//Controllers.EpisodeMenu.Attach(episodeMenuView)
	return episodeMenuView
}

func FeedMenu(controllers Controllers) *tview.List {
	feedMenu := tview.NewList()
	feedMenu.SetBorder(true)
	//Controllers.FeedMenu.Attach(feedMenu)
	return feedMenu
}

func DebugPanel(controllers Controllers) *tview.TextView {
	view := tview.NewTextView()
	return view
}
