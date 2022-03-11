package view

import "github.com/rivo/tview"

func Build() *tview.Application {
	gui := tview.NewApplication()
	gui.SetRoot(MainFlex(), true)
	return gui
}

func MainFlex() *tview.Flex {
	mainFlex := tview.NewFlex()
	mainFlex.AddItem(FeedColumn(), -1, 1, true)
	mainFlex.AddItem(EpisodeMenu(), -1, 1, true)
	return mainFlex
}

func FeedColumn() *tview.Flex {
	feedColumnFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	feedColumnFlex.AddItem(FeedMenu(), -1, 5, true)
	feedColumnFlex.AddItem(DebugPanel(), -1, 1, false)
	return feedColumnFlex
}

func EpisodeMenu() *tview.List {
	return tview.NewList()
}

func FeedMenu() *tview.List {
	feedMenu := tview.NewList()
	return feedMenu
}

func DebugPanel() *tview.TextView {
	return tview.NewTextView()
}
