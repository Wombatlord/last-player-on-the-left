package lastplayer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
)

// for Dragwing func
func drawTextLine(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for _, r := range s {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

// Contains properties for manipulating an audio stream & drawing info to the terminal. eg. volume / seeking & position
type audioPanel struct {
	buffer     *beep.Buffer
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	debug      string
}

// Constructor function for the AudioPanel struct.
func newAudioPanel(format beep.Format, streamer beep.StreamSeeker) *audioPanel {
	buffer := beep.NewBuffer(format)
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}             // used for pausing
	resampler := beep.ResampleRatio(4, 1, ctrl)                       // can change playback speed.
	volume := &effects.Volume{Streamer: streamer, Base: 2, Volume: 0} // Volume: -0.1 to 5 tested range. 0 is system volume.
	debug := ""

	return &audioPanel{
		buffer,
		format.SampleRate,
		streamer,
		ctrl,
		resampler,
		volume,
		debug,
	}
}

// Bare minimum HTTP request function to get an audio stream.
// io.ReadCloser is the interface required by mp3.Decode in StreamAudio()
// resp.Body is of this type so can simply be returned following the request.
func AudioRequest(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	return resp.Body, nil
}

// Plays the stream referenced by AudioPanel.streamer at the volume of AudioPanel.volume
func (ap *audioPanel) play() {
	speaker.Play(ap.volume)
}

// Drawing
func (ap *audioPanel) draw(screen tcell.Screen) {
	mainStyle := tcell.StyleDefault.
		Background(tcell.NewHexColor(0x473437)).
		Foreground(tcell.NewHexColor(0xD7D8A2))
	statusStyle := mainStyle.
		Foreground(tcell.NewHexColor(0xDDC074)).
		Bold(true)

	screen.Fill(' ', mainStyle)

	drawTextLine(screen, 0, 0, "This is the Last Player...ON THE LEFT!", mainStyle)
	drawTextLine(screen, 0, 1, "Press [ESC] to quit.", mainStyle)
	drawTextLine(screen, 0, 2, "Press [SPACE] to pause/resume.", mainStyle)
	drawTextLine(screen, 0, 3, "Use keys in (?/?) to turn the buttons.", mainStyle)

	speaker.Lock()
	position := ap.sampleRate.D(ap.streamer.Position())
	length := ap.sampleRate.D(ap.streamer.Len())
	volume := ap.volume.Volume
	speed := ap.resampler.Ratio()
	speaker.Unlock()

	positionStatus := fmt.Sprintf("%v / %v", position.Round(time.Second), length.Round(time.Second))
	volumeStatus := fmt.Sprintf("%.1f", volume)
	speedStatus := fmt.Sprintf("%.3fx", speed)

	drawTextLine(screen, 0, 5, "Position (Q/W):", mainStyle)
	drawTextLine(screen, 16, 5, positionStatus, statusStyle)

	drawTextLine(screen, 0, 6, "Volume   (A/S):", mainStyle)
	drawTextLine(screen, 16, 6, volumeStatus, statusStyle)

	drawTextLine(screen, 0, 7, "Speed    (Z/X):", mainStyle)
	drawTextLine(screen, 16, 7, speedStatus, statusStyle)
}

// Event handling
func (ap *audioPanel) handle(event tcell.Event) (changed, quit bool) {
	switch event := event.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyESC {
			return false, true
		}

		if event.Key() != tcell.KeyRune {
			return false, false
		}

		switch unicode.ToLower(event.Rune()) {
		case ' ':
			speaker.Lock()
			ap.ctrl.Paused = !ap.ctrl.Paused
			speaker.Unlock()
			return false, false

		case 'q', 'w':
			speaker.Lock()
			newPos := ap.streamer.Position()
			ap.debug += fmt.Sprintf("initial: %d\n", newPos)
			ap.debug += fmt.Sprintf("len: %d\n", ap.streamer.Len())
			if event.Rune() == 'q' {
				ap.debug += "q, "
				newPos -= ap.sampleRate.N(time.Second)
			}
			if event.Rune() == 'w' {
				ap.debug += "w, "
				newPos += ap.sampleRate.N(time.Second)
			}
			if newPos < 0 {
				newPos = 0
			}
			if newPos >= ap.streamer.Len() {
				newPos = ap.streamer.Len() - 1
			}

			ap.debug += fmt.Sprintf("%d, \n", newPos)

			if err := ap.streamer.Seek(newPos); err != nil {
				fmt.Print(ap.debug)
				panic(err)
			}
			speaker.Unlock()
			return true, false

		case 'a':
			speaker.Lock()
			ap.volume.Volume -= 0.1
			speaker.Unlock()
			return true, false

		case 's':
			speaker.Lock()
			ap.volume.Volume += 0.1
			speaker.Unlock()
			return true, false

		case 'z':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 15 / 16)
			speaker.Unlock()
			return true, false

		case 'x':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 16 / 15)
			speaker.Unlock()
			return true, false
		}
	}
	return false, false
}

// Proof of concept audio playback via tweaked Beep example code.
// mp3.Decode requires an io.ReadCloser
// This is provided by os.Open() for local playback
// audioRequest() provides the io.ReadCloser from a HTTP Response body for streaming from a link.
func StreamAudio(source string, audioSource string) {

	switch source {
	case "local":
		audioLocal, err := os.Open(audioSource)
		if err != nil {
			log.Fatal(err)
		}

		streamer, format, err := mp3.Decode(audioLocal)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		done := make(chan bool)
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))
		<-done

	case "stream":
		audio, err := AudioRequest(audioSource)
		if err != nil {
			log.Fatal(err)
		}

		streamer, format, err := clients.StreamDecode(audio)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		screen, err := tcell.NewScreen()
		if err != nil {
			panic(err)
		}
		err = screen.Init()
		if err != nil {
			panic(err)
		}

		ap := newAudioPanel(format, streamer)

		defer func() {
			screen.Fini()
			fmt.Println(ap.debug)
		}()

		ap.play()

		// partial prototyping to move away from beep.Seq(... func() {done <- true}) in ap.play()
		seconds := time.Tick(time.Second)
		events := make(chan tcell.Event)

		go func() {
			for {
				events <- screen.PollEvent()
			}
		}()

	loop:
		for {
			select {
			case event := <-events:
				changed, quit := ap.handle(event)
				if quit {
					break loop
				}
				if changed {
					screen.Clear()
					ap.draw(screen)
					screen.Show()
				}
			case <-seconds:
				screen.Clear()
				ap.draw(screen)
				screen.Show()
			}
		}
	}
}
