package clients

import (
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	"github.com/wombatlord/last-player-on-the-left/src/app"
)

const DownloaderLoggerName = "Downloader"

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
	logger = app.GetLogChan(DownloaderLoggerName)
	logger <- "Creating requests!"

	for _, url := range urls {
		req, err := grab.NewRequest(".", url)
		if err != nil {
			logger <- fmt.Sprintf("request error:\n%v\n%v", req, err)
		}
		reqs = append(reqs, req)
	}

	return reqs
}

func (c *DownloadClient) DownloadEpisode(client grab.Client, req *grab.Request) {
	fmt.Printf("Downloading %v... \n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)
}

func (c *DownloadClient) DownloadMulti(client grab.Client, requests ...*grab.Request) {
	for _, request := range requests {
		fmt.Printf("Downloading %v... \n", request.URL())
		client.DoBatch(-1, request)
	}
}