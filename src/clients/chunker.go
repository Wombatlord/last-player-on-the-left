package clients

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/faiface/beep"
)

const chunkSampleCount int = 1024 * 1024

type ChunkBufferStreamer struct {
	beep.StreamSeekCloser
	currentStreamer  *ClientStreamer
	nextStreamer     *ClientStreamer
	Episode          string
	totalSampleCount int
	currentOffset    int
	format beep.Format
}

func NewChunkBufferStreamer(url string) (*ChunkBufferStreamer, beep.Format) {
	// Ought to make a request for the byte range corresponding
	// to the beginning of the audio file to the chunk length
	//
	// Should record the value supplied in response header
	// that is total length of the file in bytes
	// store this value as a number of samples (we can use this to implement a len function)
	//
	// Need to know current chunkStart.
	cb := &ChunkBufferStreamer{}
	cb.Episode = url
	cb.next()
	return cb, cb.format
}

// Stream, seek, close, pos, length

func (cb *ChunkBufferStreamer) Stream(samples [][2]float64) (int, bool) {
	if cb.currentStreamer != nil {
		if cb.currentStreamer.currentStreamerExhausted() {
			n := cb.next()
			if n == nil {
				return cb.currentStreamer.Stream(samples)
			}
			cb.currentStreamer = n
		}
	} else {
		n := cb.next()
		if n == nil {
			return cb.currentStreamer.Stream(samples)
		}
		cb.currentStreamer = n
	}

	return cb.currentStreamer.Stream(samples)
}

func (cb *ChunkBufferStreamer) Seek(p int) error {
	beforeCurrentOffset := p < cb.currentOffset 
	afterNextEnd := p > cb.currentOffset + chunkSampleCount * 2
	
	if beforeCurrentOffset || afterNextEnd {
		cb.currentStreamer = nil
		cb.nextStreamer = nil
		cb.currentOffset = p
		cb.currentStreamer = cb.next()
		return cb.currentStreamer.Seek(0)
	}

	beforeNextEnd := p < cb.currentOffset + chunkSampleCount * 2
	afterNextStart := p >= cb.currentOffset + chunkSampleCount
	
	if afterNextStart && beforeNextEnd {
		cb.currentStreamer = cb.next()
		return cb.currentStreamer.Seek(p - cb.currentOffset)
	}
	
	return cb.currentStreamer.Seek(p - cb.currentOffset)
}

func (cb *ChunkBufferStreamer) Close() error {
	if cb.currentStreamer != nil {
		_ = cb.currentStreamer.Close()
	}

	if cb.nextStreamer != nil {
		_ = cb.nextStreamer.Close()
	}

	return nil
}

func (cb *ChunkBufferStreamer) Position() int {
	return cb.currentStreamer.Position() + cb.currentOffset
}

func (cb *ChunkBufferStreamer) Length() int {
	return cb.totalSampleCount
}

func (cb *ChunkBufferStreamer) RequestChunkAtOffset(offset int) io.ReadCloser {
	rc, _ := cb.RequestSampleRange(offset, offset+chunkSampleCount)
	return rc
}

func (cb *ChunkBufferStreamer) RequestSampleRange(start, end int) (io.ReadCloser, *http.Response) {
	startByte := 4 * start
	endByte := 4 * end

	rangeHeader := fmt.Sprintf("bytes=%d-%d", startByte, endByte-1)
	req, _ := http.NewRequest("GET", cb.Episode, nil)
	req.Header.Set("Range", rangeHeader)
	resp, _ := http.DefaultClient.Do(req)

	if cb.totalSampleCount == 0 {
		contentRange := resp.Header.Get("Content-Range")
		totalByteLength, _ := strconv.Atoi(strings.Split(contentRange, "/")[1])
		cb.totalSampleCount = totalByteLength / 4
	}

	return resp.Body, resp
}

func (cb *ChunkBufferStreamer) next() *ClientStreamer {
	nextOffset := cb.currentOffset + chunkSampleCount
	format := beep.Format{}
	streamer := &ClientStreamer{}
	nextStreamer := &ClientStreamer{}

	if cb.currentStreamer == nil && cb.nextStreamer == nil {
		// this is the first streamer starting at the offset in the current streamer
		// and the second streamer starting at the offset in the next streamer.
		streamer, format, _ = StreamDecode(cb.RequestChunkAtOffset(cb.currentOffset))
		cb.currentStreamer = streamer
		
		nextStreamer, format, _ = StreamDecode(cb.RequestChunkAtOffset(nextOffset))
		cb.nextStreamer = nextStreamer
	}

	if cb.currentStreamer != nil && cb.nextStreamer != nil {
		// both streamers are populated, so copy nextStreamer to currentStreamer
		// and initiate the next streamer.
		cb.currentStreamer = cb.nextStreamer

		if nextOffset < cb.totalSampleCount {
			streamer, format, _ = StreamDecode(cb.RequestChunkAtOffset(nextOffset))
			cb.nextStreamer = streamer
		} else {
			cb.nextStreamer = nil
		}
	}

	cb.currentOffset = nextOffset
	emptyFormat := beep.Format{}
	if format == emptyFormat {
		cb.format = format
	}

	// These cases should never happen.
	if cb.currentStreamer == nil && cb.nextStreamer != nil {
		return nil
	}

	if cb.currentStreamer != nil && cb.nextStreamer == nil {
		return nil
	}

	return cb.currentStreamer
}
