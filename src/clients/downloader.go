package clients

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"os"
	"time"
)

type DownloadClient struct {
	Client *grab.Client
}

func NewClient() *grab.Client {
	// create a client
	client := grab.NewClient()

	// set the User Agent header
	client.UserAgent = "Last Player On The Left"

	return client
}

func (c *DownloadClient) CreateRequests(urls []string) (reqs []*grab.Request) {
	loggers[DLLog].Print("Creating requests!")

	for _, url := range urls {
		req, err := grab.NewRequest(".", url)
		if err != nil {
			loggers[DLLog].Printf("request error:\n%v\n%v", req, err)
		}
		reqs = append(reqs, req)
	}

	return reqs
}

func (c *DownloadClient) DownloadEpisode(client grab.Client, req *grab.Request) {
	fmt.Printf("Downloading %v... \n", req.URL())
	resp := c.Client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	// Output:
	// Downloading http://www.golang-book.com/public/pdf/gobook.pdf...
	//   200 OK
	//   transferred 42970 / 2893557 bytes (1.49%)
	//   transferred 1207474 / 2893557 bytes (41.73%)
	//   transferred 2758210 / 2893557 bytes (95.32%)
	// Download saved to ./gobook.pdf
}

func (c *DownloadClient) DownloadMulti(client grab.Client, requests ...*grab.Request) {
	responses := c.Client.DoBatch(-1, requests...)

	go func() {
		for elem := range responses {
			loggers[DLLog].Print(elem.HTTPResponse.Request.Host)
			loggers[DLLog].Print(elem.Request.HTTPRequest.UserAgent())
			loggers[DLLog].Print(elem.Filename)
			loggers[DLLog].Print(elem.HTTPResponse.Status)
		}
	}()

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
}
