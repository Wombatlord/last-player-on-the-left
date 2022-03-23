package clients

import (
	"fmt"
	"github.com/faiface/beep"
	gomp3 "github.com/hajimehoshi/go-mp3"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Success is the named constructor for the case that the download completes
// successfully
func Success(size int64) *SizedResult {
	return &SizedResult{Size: size, Err: nil}
}

// Fail is the named constructor for the case that the copy to disk fails
func Fail(failure error) *SizedResult {
	return &SizedResult{Size: 0, Err: failure}
}

// SizedResult is the struct that encapsulates the error on failure or the
// size, either predicted or written to disk, of the corresponding data in bytes
type SizedResult struct {
	Size int64
	Err  error
}

// IsSuccess returns true if the SizedResult.Err is not nil, i.e. a return
// value of true corresponds to the file being successfully downloaded
func (r *SizedResult) IsSuccess() bool {
	return r.Err == nil
}

const (
	gomp3NumChannels   = 2
	gomp3Precision     = 2
	gomp3BytesPerFrame = gomp3NumChannels * gomp3Precision
)

// boopDecode takes a ReadCloser containing audio data in MP3 format and returns a StreamSeekCloser,
// which streams that audio. The Seek method will panic if rc is not io.Seeker.
//
// Do not close the supplied ReadSeekCloser, instead, use the Close method of the returned
// StreamSeekCloser when you want to release the resources.
func boopDecode(rc io.ReadCloser) (s *Decoder, format beep.Format, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "mp3")
		}
	}()
	d, err := gomp3.NewDecoder(rc)
	if err != nil {
		return nil, beep.Format{}, err
	}
	format = beep.Format{
		SampleRate:  beep.SampleRate(d.SampleRate()),
		NumChannels: gomp3NumChannels,
		Precision:   gomp3Precision,
	}
	return &Decoder{rc, d, format, 0, nil, 0}, format, nil
}

// Decoder is a ripoff of the streamer that mp3.Decode returns but with a Decoder.SetLength. All the methods
// are identical except the Decoder.SetLength.
type Decoder struct {
	closer io.Closer
	d      *gomp3.Decoder
	f      beep.Format
	pos    int
	err    error
	len    int
}

func (d *Decoder) Stream(samples [][2]float64) (n int, ok bool) {
	if d.err != nil {
		return 0, false
	}
	var tmp [gomp3BytesPerFrame]byte
	for i := range samples {
		dn, err := d.d.Read(tmp[:])
		if dn == len(tmp) {
			samples[i], _ = d.f.DecodeSigned(tmp[:])
			d.pos += dn
			n++
			ok = true
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			d.err = errors.Wrap(err, "mp3")
			break
		}
	}
	return n, ok
}

func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) Len() int {
	return d.len / gomp3BytesPerFrame
}

func (d *Decoder) SetLength(length int) {
	d.len = length
}

func (d *Decoder) Position() int {
	return d.pos / gomp3BytesPerFrame
}

func (d *Decoder) Seek(p int) error {
	if p < 0 || d.Len() < p {
		return fmt.Errorf("mp3: seek position %v out of range [%v, %v]", p, 0, d.Len())
	}
	_, err := d.d.Seek(int64(p)*gomp3BytesPerFrame, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "mp3")
	}
	d.pos = p * gomp3BytesPerFrame
	return nil
}

func (d *Decoder) Close() error {
	err := d.closer.Close()
	if err != nil {
		return errors.Wrap(err, "mp3")
	}
	return nil
}

// TcpDiskBufferedStreamer is basically a proxy to some hacked together beep source code, all credit to them, unless
// the code looks a bit messy, that's probably us.
func TcpDiskBufferedStreamer(url string) (streamer *Decoder, format beep.Format) {
	filepath := "cache/cache.mp3"
	var (
		done    chan *SizedResult
		started chan *SizedResult
	)

	if fileInfo, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		started, done = asyncDownloadAudio(filepath, url)
	} else {
		go func() {
			started <- Success(fileInfo.Size())
			done <- Success(fileInfo.Size())
		}()
	}

	streamer, format = getStreamer(started, filepath)

	go func() {
		if result := <-done; result.IsSuccess() {
			streamer.SetLength(int(result.Size))
		}
	}()

	return streamer, format
}

// getStreamer uses the filesystem path of the cached audio to create the Decoder. It mirrors the interface of
// mp3.Decode from beep
func getStreamer(started chan *SizedResult, filepath string) (streamer *Decoder, format beep.Format) {
	if result := <-started; result.IsSuccess() {
		audio, err := os.Open(filepath)
		if err != nil {
			log.Fatalln(err)
		}
		streamer, format, err = boopDecode(audio)
		if err != nil {
			log.Fatalln(err)
		}
		streamer.SetLength(int(result.Size))
	} else {
		log.Fatalln(result)
	}
	return streamer, format
}

// asyncDownloadAudio sets off a doDownload goroutine and returns the started and done channels that
// it will report back its progress on.
func asyncDownloadAudio(filename, url string) (started chan *SizedResult, done chan *SizedResult) {
	started, done = make(chan *SizedResult), make(chan *SizedResult)

	go doDownload(filename, url, started, done)

	return started, done
}

// doDownload will download the provided url to the filepath specified. Supply two chan error, started
// will send nil on successful start, or an error. If started returns nil, then done will send another
// nil if the file successfully downloaded and the error otherwise.
func doDownload(filepath string, url string, started chan *SizedResult, done chan *SizedResult) {
	defer close(started)

	var out *os.File

	// Get the data
	resp, err := http.Get(url)
	fmt.Printf("Content Length: %d\n", resp.ContentLength)
	if err != nil {
		started <- Fail(err)
	}

	// Create the file
	out, err = os.Create(filepath)
	if err != nil {
		started <- Fail(err)
	}

	go func() {
		defer func() {
			close(done)
		}()
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		defer func(out *os.File) {
			// might need to handle this if it's an issue
			_ = out.Close()
		}(out)

		// Get to Copyin'
		size, copyErr := io.Copy(out, resp.Body)

		// Send the result down the pipe
		if copyErr != nil {
			fmt.Println(copyErr)
			done <- Fail(copyErr)
		} else {
			done <- Success(size)
		}
	}()

	// Wait a little for the copy to disk to begin
	time.Sleep(time.Second / 10)
	// then signal the download has begun
	started <- Success(resp.ContentLength)
}
