package lastplayer

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

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

		done := make(chan bool)
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))
		<-done
	}
}
