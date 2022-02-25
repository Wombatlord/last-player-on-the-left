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

// Contains properties for manipulating an audio stream & drawing info to the terminal. eg. volume / seeking & position
type audioPanel struct {
	sampleRate beep.SampleRate
	streamer beep.Streamer
	ctrl *beep.Ctrl
	resampler *beep.Resampler
	volume *effects.Volume
}

// Constructor function for the AudioPanel struct.
func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *audioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)} // used for pausing
	resampler := beep.ResampleRatio(4, 1, ctrl)	// can change playback speed.
	volume := &effects.Volume{Streamer: streamer, Base: 2, Volume: 0} // Volume: -0.1 to 5 tested range. 0 is system volume.
	return &audioPanel{sampleRate, streamer, ctrl, resampler, volume}
}

// Plays the stream referenced by AudioPanel.streamer at the volume of AudioPanel.volume
func (ap *audioPanel) play() {
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

		ap := newAudioPanel(format.SampleRate, streamer)
		
		ap.play()
		seconds := time.Tick(time.Second)
		
		for {
			select {
			case <- seconds:
				fmt.Println("playing")
			}
		}
	}
}