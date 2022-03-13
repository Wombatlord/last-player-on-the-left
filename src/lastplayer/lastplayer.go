package lastplayer

import (
	"fmt"
	"github.com/rivo/tview"
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
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
)

// drawTextLine is a convenience function for drawing text
func drawTextLine(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for _, r := range s {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

var (
	panel = &AudioPanel{}
)

func getLogger() chan string {
	return app.GetLogChan("lastplayer")
}

// PlayerState is a snapshot of the current player state
type PlayerState struct {
	Position time.Duration
	Length   time.Duration
	Playing  bool
}

type PlayerStateSubscriber interface {
	OnUpdate()
}

// AudioPanel contains properties for manipulating an audio stream & drawing info to the terminal. eg. volume / seeking & position
type AudioPanel struct {
	sampleRate  beep.SampleRate
	streamer    beep.StreamSeeker
	ctrl        *beep.Ctrl
	resampler   *beep.Resampler
	volume      *effects.Volume
	subscribers []PlayerStateSubscriber
	gui         *tview.Application
	ticker      *time.Ticker
}

func (ap *AudioPanel) PlayPause() {
	speaker.Lock()
	ap.ctrl.Paused = !ap.ctrl.Paused
	speaker.Unlock()
}

// InitAudioPanel is a constructor function for the AudioPanel struct.
func InitAudioPanel(format beep.Format, streamer beep.StreamSeeker) *AudioPanel {
	getLogger() <- "Building audio panel"
	panel.subscribers = []PlayerStateSubscriber{}
	panel.SetStreamer(format, streamer)
	return panel
}

// FetchAudioPanel will return the already initialised panel pointer.
func FetchAudioPanel() *AudioPanel {
	return panel
}

func (ap *AudioPanel) SubscribeToState(subscriber PlayerStateSubscriber) {
	ap.subscribers = append(ap.subscribers, subscriber)
}

func (ap *AudioPanel) AttachApp(gui *tview.Application) {
	ap.gui = gui
}

func (ap *AudioPanel) SetStreamer(format beep.Format, streamer beep.StreamSeeker) {
	ap.streamer = streamer
	ap.sampleRate = format.SampleRate
	ap.ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}            // used for pausing
	ap.resampler = beep.ResampleRatio(4, 1, ap.ctrl)                   // can change playback speed.
	ap.volume = &effects.Volume{Streamer: ap.ctrl, Base: 2, Volume: 0} // Volume: -0.1 to 5 tested range. 0 is system volume
}

func (ap *AudioPanel) SpawnPublisher() {
	ap.ticker = time.NewTicker(time.Second)
	publisher := func() {
		for range ap.ticker.C {
			for _, sub := range ap.subscribers {
				ap.gui.QueueUpdateDraw(func() {
					sub.OnUpdate()
				})
			}
		}
	}
	go publisher()
}

func (ap *AudioPanel) PlayFromUrl(url string) {
	var err error
	getLogger() <- "PlayFromUrl call"

	audio, err := AudioRequest(url)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := clients.StreamDecode(audio)
	if err != nil {
		log.Fatal(err)
	}
	ap.SetStreamer(format, streamer)

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}
	ap.play()
}

// AudioRequest is the bare minimum HTTP request function to get an audio stream.
// io.ReadCloser is the interface required by mp3.Decode in StreamAudio()
// resp.Body is of this type so can simply be returned following the request.
func AudioRequest(url string) (io.ReadCloser, error) {
	getLogger() <- fmt.Sprintf("Requesting streaming audio from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	return resp.Body, nil
}

// play Plays the stream referenced by AudioPanel.streamer at the volume of AudioPanel.volume
func (ap *AudioPanel) play() {
	getLogger() <- "Playing audio"
	speaker.Play(ap.volume)
}

// GetPlayerState returns a PlayerState value that represents a snapshot of the user relevant
// player state at the time this method was called
func (ap *AudioPanel) GetPlayerState() PlayerState {
	state := PlayerState{}
	if ap.streamer != nil {
		if ap.streamer.Position() <= 0 {
			return state
		}
		speaker.Lock()
		state.Position = ap.sampleRate.D(ap.streamer.Position())
		state.Length = ap.sampleRate.D(ap.streamer.Len())
		state.Playing = !ap.ctrl.Paused
		speaker.Unlock()
	}

	return state
}

// Drawing
func (ap *AudioPanel) draw(screen tcell.Screen) {
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

// AppEvent handling
func (ap *AudioPanel) handle(eventInstance tcell.Event) (changed, quit bool) {
	switch event := eventInstance.(type) {
	case *tcell.EventKey:
		getLogger() <- fmt.Sprintf("Handling AppEvent. Type: %T, Name: %s", eventInstance, event.Name())
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
			if event.Rune() == 'q' {
				newPos -= ap.sampleRate.N(time.Second)
			}
			if event.Rune() == 'w' {
				newPos += ap.sampleRate.N(time.Second)
			}
			if newPos < 0 {
				newPos = 0
			}
			if newPos >= ap.streamer.Len() {
				newPos = ap.streamer.Len() - 1
			}

			if err := ap.streamer.Seek(newPos); err != nil {
				log.Fatal(err)
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

// StreamAudio is a proof of concept audio playback via tweaked Beep example code.
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
		defer func(streamer beep.StreamSeekCloser) {
			err := streamer.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(streamer)

		_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

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
		defer func(streamer *clients.ClientStreamer) {
			err := streamer.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(streamer)

		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			log.Fatal(err)
		}

		screen, err := tcell.NewScreen()
		if err != nil {
			panic(err)
		}
		err = screen.Init()
		if err != nil {
			panic(err)
		}

		ap := InitAudioPanel(format, streamer)

		defer func() {
			screen.Fini()
		}()

		ap.play()

		// partial prototyping to move away from beep.Seq(... func() {done <- true}) in ap.play()
		events := make(chan tcell.Event)

		go func() {
			for {
				events <- screen.PollEvent()
			}
		}()

		t := time.NewTicker(time.Second)
		defer t.Stop()

	loop:
		for {
			select {
			case event := <-events:
				changed, quit := ap.handle(event)
				if quit {
					getLogger() <- "Exiting"
					break loop
				}
				if changed {
					screen.Clear()
					ap.draw(screen)
					screen.Show()
				}
			case <-t.C:
				screen.Clear()
				ap.draw(screen)
				screen.Show()
			}
		}
	}
}
