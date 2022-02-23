package lastplayer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	// "github.com/gdamore/tcell"
)

type AudioPanel struct {
	sampleRate beep.SampleRate
	streamer beep.Streamer
	ctrl *beep.Ctrl
	resampler *beep.Resampler
	volume *effects.Volume
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *AudioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &AudioPanel{sampleRate, streamer, ctrl, resampler, volume}
}

func (ap *AudioPanel) play() {
	speaker.Play(ap.volume)
}

// Bare minimum HTTP request function to get an audio stream.
// io.ReadCloser is the interface required by mp3.Decode in StreamAudio()
// resp.Body is of this type so can simply be returned following the request.
func audioRequest(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	
	return resp.Body, nil
}

// Proof of concept audio playback via tweaked Beep example code. 
// mp3.Decode requires an io.ReadCloser
// This is provided by os.Open() for local playback
// audioRequest() provides the io.ReadCloser from a HTTP Response body for streaming from a link.
func StreamAudio(source string, audioSource string) {

	switch source{
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
		audio, err := audioRequest(audioSource)
		if err != nil {
			log.Fatal(err)
		}

		streamer, format, err := mp3.Decode(audio)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
		speaker.Play(ctrl)

		for {
			fmt.Print("Press [ENTER] to pause/resume. ")
			fmt.Scanln()

			speaker.Lock()
			ctrl.Paused = !ctrl.Paused
			speaker.Unlock()
		}
	}
}