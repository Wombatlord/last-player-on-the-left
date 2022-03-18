package clients

import (
	"encoding/json"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"io"
)

// StreamDecode is replicating the interface of mp3.Decode returning a buffer instead of a reader.
// It acts as a wrapper to asynchronously buffer the tcp audio stream.
func StreamDecode(audio io.ReadCloser) (*ClientStreamer, beep.Format, error) {
	loggers[TCPALog].Print("Decoding TCP Audio Stream")
	var (
		streamer beep.StreamCloser
		format   beep.Format
		err      error
	)
	streamer, format, err = mp3.Decode(audio)
	if err != nil {
		panic(err)
	}
	buff := beep.NewBuffer(format)

	clientStreamer := &ClientStreamer{Buff: buff, inputStreamer: streamer, quit: false}

	go func() {
		buff.Append(streamer)
		clientStreamer.quit = true
		loggers[TCPALog].Print("Buffering goroutine complete")
	}()

	return clientStreamer, format, err
}

type ClientStreamer struct {
	Buff            *beep.Buffer
	inputStreamer   beep.StreamCloser
	quit            bool
	currentStreamer beep.StreamSeeker
}

type streamerState struct {
	BufferedSamples  int
	StreamingSamples int
	Position         int
	Quit             bool
	StreamerErr      error
}

func (client *ClientStreamer) String() string {
	info := &streamerState{
		BufferedSamples:  client.Buff.Len(),
		StreamingSamples: client.currentStreamer.Len(),
		Position:         client.currentStreamer.Position(),
		Quit:             client.quit,
		StreamerErr:      client.Err(),
	}

	jsonData, err := json.Marshal(info)
	if err != nil {
		jsonData = []byte(err.Error())
	}

	return string(jsonData)
}

func (client *ClientStreamer) currentStreamerExhausted() bool {
	return client.currentStreamer.Position() == client.currentStreamer.Len()
}

func (client *ClientStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	// if not first iteration
	if client.currentStreamer != nil {
		if client.currentStreamerExhausted() {
			client.currentStreamer = client.next()
		}
	} else {
		client.currentStreamer = client.next()
	}

	return client.currentStreamer.Stream(samples)
}

func (client *ClientStreamer) next() beep.StreamSeeker {
	streamer := client.Buff.Streamer(0, client.Buff.Len())

	if client.currentStreamer != nil {
		err := streamer.Seek(client.currentStreamer.Position())
		if err != nil {
			return streamer
		}
	}

	return streamer
}

func (client *ClientStreamer) Close() error {
	return client.inputStreamer.Close()
}

func (client *ClientStreamer) Err() error {
	if client.currentStreamer == nil {
		return nil
	}

	return client.currentStreamer.Err()
}

func (client *ClientStreamer) Len() int {
	return client.Buff.Len()
}

func (client *ClientStreamer) Position() int {
	if client.currentStreamer == nil {
		return 0
	}
	return client.currentStreamer.Position()
}

func (client *ClientStreamer) Seek(p int) error {
	client.currentStreamer = client.Buff.Streamer(0, client.Buff.Len())
	return client.currentStreamer.Seek(p)
}
