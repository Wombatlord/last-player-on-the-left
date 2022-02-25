package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

type Alignment bool

const (
	Vertical   Alignment = true
	Horizontal           = false
)

type Box struct {
	Percent int
	Split   *[2]Box
	Align   Alignment
}

// ui is the display laid out like this:
//
//	+-------------------------------+
//	|				|				|
//	|				|				|
//	|				|				|
//	|				|				|
//	+-------------------------------+
//	| 								|
//	+-------------------------------+
//
var ui = Box{
	Percent: 100,
	Align:   Horizontal,
	Split: &[2]Box{
		Box{
			Percent: 80,
			Align:   Vertical,
			Split: &[2]Box{
				Box{
					Percent: 50,
				},
				Box{
					Percent: 50,
				},
			},
		},
		Box{
			Percent: 20,
		},
	},
}

var (
	w       int
	h       int
	selecta int
)

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func RuneLines() [][]rune {
	w, err := os.ReadFile("lorem.txt")
	if err != nil {
		panic(err)
	}
	lines := bytes.Split(w, []byte("\n"))
	runeLines := make([][]rune, len(lines))
	for i, line := range lines {
		runeLines[i] = []rune(string(line))
	}
	return runeLines
}

func words() string {
	w, err := os.ReadFile("lorem.txt")
	if err != nil {
		panic(err)
	}
	return string(w)
}

func main() {
	selecta = 2
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	w, h = s.Size()

	// Draw initial boxes
	drawBox(s, 30, 15, 57, 20, boxStyle, "Press C to reset")
	menu(s, boxStyle)

	// Event loop
	ox, oy := -1, -1
	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	for {
		w, h = s.Size()
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		drawText(s, w/2, (h/2)-1, w-1, h/2, defStyle, fmt.Sprintf("%5d", selecta))

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}

			switch ev.Key() {
			case tcell.KeyUp:
				selecta--
				selecta = max(selecta, 1)
				drawSelectorHighlight(s, 1, selecta, w/3-1, boxStyle)

				s.Show()
			case tcell.KeyDown:
				selecta++
				selecta = min(selecta, h-1)
				drawSelectorHighlight(s, 1, selecta, w/3-1, boxStyle)

				s.Show()
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			// Only process button events, not wheel events
			button &= tcell.ButtonMask(0xff)

			// Record on
			if button != tcell.ButtonNone && ox < 0 {
				ox, oy = x, y
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if ox >= 0 {
					label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
					drawBox(s, ox, oy, x, y, boxStyle, label)
					ox, oy = -1, -1
				}
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func menu(s tcell.Screen, boxStyle tcell.Style) {
	top := 0
	bottom := h - 1
	drawBox(s, 0, top, w/3, bottom, boxStyle, "WORDS")
	for i, line := range RuneLines() {
		if i+top >= bottom-1 {
			break
		}
		sLine := string(line)
		drawText(s, 1, i+1, w/3-1, i+2, boxStyle, sLine)
	}
}

func drawSelectorHighlight(s tcell.Screen, x, y, width int, defaultStyle tcell.Style) {
	for i := 0; i < width; i++ {
		col := i + x
		row := y
		mainc, _, style, _ := s.GetContent(col, row)
		highlight := style.Background(tcell.ColorDarkOrange)
		s.SetContent(col, row, mainc, nil, highlight)
	}
	for _, j := range []int{-1, 1} {
		for i := 0; i < width; i++ {
			col := i + x
			row := y + j
			mainc, _, _, _ := s.GetContent(col, row)
			s.SetContent(col, row, mainc, nil, defaultStyle)
		}
	}
}
