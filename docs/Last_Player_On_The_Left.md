# Last Player On The Left

## Properties
**Languange:** Go
**Interface:** Terminal UI

## Features
- RSS Feed Subscription
- Audio streaming from links in parsed RSS XML.
- Download episodes for offline playback.
- Terminal Interface for RSS subscription, browsing, playback.

## Design
"The Last Player On The Left" (LPOTL) is an in terminal audio streaming app. The goal is to be simple, lightweight, and without external dependencies on programs like VLC or MPV.


##### RSS Subscription
LPOTL retrieves XML data from RSS feeds via HTTP requests using Golangs "net/html" library. Using **xml.Unmarshal** to parse this data into structs allows both programmatic access via the xml tags, alongside the ability to serialise the data locally.

Direct links to episode streams can be retrieved from either the Go struct, or from the local database. 

- Request & Parse XML [Done!]
- Serialise XML to database [Proof of Concept!]
	- currently serialising to `config.yaml`, see `config.go`.


##### Audio Streaming
Links to audio streams are enclosed by the **link** tag in an RSS feed. At the top level, the link tag likely references the domain of a given podcast.

Episode links are found in the **link** tag contained within **item** tags. Each **item** refers to a single episode. This is the link that should be retrieved and passed to an interface, either for streaming or download and playback.

- **Find an appropriate audio library [Done!]
	- **beep https://github.com/faiface/beep
		- Beep has been chosen! Easy, lightweight, and does not require building seperate dependancies like portaudio.

- **Implement an interface to stream audio from a link [Done!]
	- `tcp_audio.go` contains an implementation of `mp3.Decode()` in `StreamDecode()`, returning a buffer instead of a reader. `StreamDecode()` acts as a wrapper around `mp3.Decode()` to asynchronously buffer the tcp audio stream.
	- `StreamDecode()` returns a `ClientStreamer` struct, which contains:
		- `beep.Buffer` which is streamed into via the `inputStreamer`
		- `inputStreamer` is the initial audio stream decoding from `mp3.Decode`
		- `quit` bool for indicating the end of a stream or otherwise exhausted buffer.
		- `currentStreamer` which is a `beep.StreamSeeker` used to seek through the `ClientStreamer.Buff`
	- `ClientStreamer` implements the interfaces of beep.Streamer as attached methods which proxy back to the beep.Streamer interfaces. This allows `StreamDecode` to be a simple wrapper around `mp3.Decode` with the returned `ClientStreamer` retaining access to the original `beep.Streamer` interfaces.

- **Implement an interface to play audio from a local file [Proof Of Concept!]
	- Beep's `mp3.Decode()` is able to handle both local playback & streaming through the `io.ReadCloser` interface.


- **Download an audio file from the RSS Feed for local playback [Pending]

#### Logging

Last Player provides application logs in logs/log.txt. The log origin is indicated by the content contained in square brackets, which identifies the package from which the log originated. Pass the log file to grep to search for specific log  entires.

eg: `$cat log.txt | grep "main"`

##### Terminal UI
Provide a simple terminal interface for searching for podcasts, subscribing to the feed, and streaming / downloading episodes.

If a local database exists, LPOTL should retrieve data to populate the interface with subscribed content.

- "github.com/gdamore/tcell" [Chosen!]