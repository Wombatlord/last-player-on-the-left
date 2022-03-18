package audiopanel

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	panel = &AudioPanel{}
)

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
	clock       *time.Ticker
	Format      beep.Format
	logger      *log.Logger
	callback    func(func())
}

func (ap *AudioPanel) PlayPause() {
	if ap.ctrl == nil {
		return
	}
	speaker.Lock()
	ap.ctrl.Paused = !ap.ctrl.Paused
	speaker.Unlock()
}

func (ap *AudioPanel) Duration(e clients.Enclosure) time.Duration {
	byteCount := int(e.Length)
	numSamples := byteCount / ap.Format.Width()
	return ap.sampleRate.D(numSamples)
}

// FetchAudioPanel will return the already initialised panel pointer.
func FetchAudioPanel() *AudioPanel {
	return panel
}

func (ap *AudioPanel) SubscribeToState(subscriber PlayerStateSubscriber) {
	ap.subscribers = append(ap.subscribers, subscriber)
}

// AttachLogger implements setter injection of a log.Logger
func (ap *AudioPanel) AttachLogger(logger *log.Logger) *AudioPanel {
	ap.logger = logger
	return ap
}

// SetPublishCallback is where the function to update on publish should be supplied
// Probably don't pass anything other than QueueUpdateDraw
func (ap *AudioPanel) SetPublishCallback(callback func(func())) {
	ap.callback = callback
}

func (ap *AudioPanel) SetStreamer(format beep.Format, streamer beep.StreamSeeker) {
	ap.Format = format
	ap.streamer = streamer
	ap.sampleRate = format.SampleRate
	ap.ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}            // used for pausing
	ap.resampler = beep.ResampleRatio(4, 1, ap.ctrl)                   // can change playback speed.
	ap.volume = &effects.Volume{Streamer: ap.ctrl, Base: 2, Volume: 0} // Volume: -0.1 to 5 tested range. 0 is system volume
}

func (ap *AudioPanel) SpawnPublisher() {
	ap.clock = time.NewTicker(time.Second / 2)
	publisher := func() {
		for range ap.clock.C {
			for _, sub := range ap.subscribers {
				ap.callback(func() {
					sub.OnUpdate()
				})
			}
		}
	}
	go publisher()
}

func (ap *AudioPanel) PlayFromUrl(url string) {
	var err error
	ap.logger.Println("PlayFromUrl call")

	audio, err := ap.AudioRequest(url)
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
func (ap *AudioPanel) AudioRequest(url string) (io.ReadCloser, error) {
	ap.logger.Printf("Requesting streaming audio from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	return resp.Body, nil
}

// play Plays the stream referenced by AudioPanel.streamer at the volume of AudioPanel.volume
func (ap *AudioPanel) play() {
	ap.logger.Println("Playing audio")
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
