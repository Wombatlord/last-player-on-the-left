package view

import (
	"github.com/rivo/tview"
)

type Controllers struct {
	FeedMenu         *FeedsMenuController
	EpisodeMenu      *EpisodeMenuController
	RootController   *RootController
	APViewController *APViewController
}

// Build returns the tview app, implement any additions to
// the user interface adding new primitives to the app hierarchy
// in this file
func Build(gui *tview.Application, controllers Controllers) *tview.Application {
	mainFlex := MainFlex(controllers).SetDirection(tview.FlexRow)
	topRow := TopRow()
	episodeMenu := EpisodeMenu(controllers)
	feedMenu := FeedMenu(controllers)
	apView := AudioPanelView(controllers)

	topRow.AddItem(feedMenu, -1, 1, true)
	topRow.AddItem(episodeMenu, -1, 1, true)

	mainFlex.AddItem(topRow, -1, 4, true)
	mainFlex.AddItem(apView, -1, 1, false)

	gui.SetRoot(mainFlex, true)

	focusRing := []tview.Primitive{feedMenu, episodeMenu}
	controllers.RootController.SetFocusRing(focusRing)

	return gui
}

func MainFlex(controllers Controllers) *tview.Flex {
	mainFlex := tview.NewFlex()
	controllers.RootController.Attach(mainFlex)

	return mainFlex
}

func TopRow() *tview.Flex {
	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	return topRow
}

func EpisodeMenu(controllers Controllers) *tview.List {
	episodeMenuView := tview.NewList()
	controllers.EpisodeMenu.Attach(episodeMenuView)

	episodeMenuView.SetBorder(true).
		SetTitle("Episodes").
		SetTitleAlign(tview.AlignCenter)
	return episodeMenuView
}

func FeedMenu(controllers Controllers) *tview.List {
	feedMenu := tview.NewList()
	controllers.FeedMenu.Attach(feedMenu)

	feedMenu.SetBorder(true).
		SetTitle("Podcasts").
		SetTitleAlign(tview.AlignCenter)
	return feedMenu
}

func AudioPanelView(controllers Controllers) *tview.TextView {
	view := tview.NewTextView()
	controllers.APViewController.Attach(view)

	view.SetBorder(true).
		SetTitle("Player").
		SetTitleAlign(tview.AlignCenter)
	return view
}
