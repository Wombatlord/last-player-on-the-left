package view

import (
	"github.com/rivo/tview"
)

// This file is for convenience functions that
// initialise the various views that make up the
// LastPlayer. This is where to configure things
// like the border/title etc.

func MainFlex() *tview.Flex {
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	return mainFlex
}

func TopRow() *tview.Flex {
	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	return topRow
}

func EpisodeMenu() *tview.List {
	episodeMenuView := tview.NewList()

	episodeMenuView.SetBorder(true).
		SetTitle("Episodes").
		SetTitleAlign(tview.AlignCenter)
	return episodeMenuView
}

func FeedMenu() *tview.List {
	feedMenu := tview.NewList()

	feedMenu.SetBorder(true).
		SetTitle("Podcasts").
		SetTitleAlign(tview.AlignCenter)
	return feedMenu
}

func AudioPanelView() *tview.TextView {
	view := tview.NewTextView()

	view.SetBorder(true).
		SetTitle("Player").
		SetTitleAlign(tview.AlignCenter)
	return view
}
