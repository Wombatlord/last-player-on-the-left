package clients

import (
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	"github.com/wombatlord/last-player-on-the-left/src/app"
)

const DownloaderLoggerName = "Downloader"

var dlLogger chan string

type DownloadClient struct {
	Client *grab.Client
}

func NewClient() *grab.Client {
	dlLogger = app.GetLogChan(DownloaderLoggerName)

	// create a client
	client := grab.NewClient()

	// set the User Agent header
	client.UserAgent = "Last Player On The Left"

	return client
}

func (c *DownloadClient) CreateRequests(urls []string) (reqs []*grab.Request) {
	dlLogger = app.GetLogChan(DownloaderLoggerName)
	dlLogger <- "Creating requests!"

	for _, url := range urls {
		req, err := grab.NewRequest(".", url)
		if err != nil {
			dlLogger <- fmt.Sprintf("request error:\n%v\n%v", req, err)
		}
		reqs = append(reqs, req)
	}

	return reqs
}

func (c *DownloadClient) DownloadEpisode(client grab.Client, req *grab.Request) {
	fmt.Printf("Downloading %v... \n", req.URL())
	resp := c.Client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)
}

func (c *DownloadClient) DownloadMulti(client grab.Client, requests ...*grab.Request) {
	dlLogger = app.GetLogChan(DownloaderLoggerName)
	responses := c.Client.DoBatch(-1, requests...)

	go func() {
		for elem := range responses {
			dlLogger <- elem.HTTPResponse.Request.Host
			dlLogger <- elem.Request.HTTPRequest.UserAgent()
			dlLogger <- elem.Filename
			dlLogger <- elem.HTTPResponse.Status
		}
	}()

	// Only one log message currently being received
	// dlLogger <- responses.Err().Error()
	// dlLogger <- responses.HTTPResponse.Request.Host
	// dlLogger <- responses.Request.HTTPRequest.UserAgent()
	// dlLogger <- responses.Filename
	// dlLogger <- responses.HTTPResponse.Status

	// for _, request := range requests {
	// 	fmt.Printf("Downloading %v... \n", request.URL())
	// 	c.Client.DoBatch(-1, request)
	// 	dlLogger <- request.HTTPRequest.Host
	// 	dlLogger <- request.HTTPRequest.Method
	// 	dlLogger <- c.Client.UserAgent
	// 	dlLogger <- request.HTTPRequest.UserAgent()
	// 	dlLogger <- request.Filename
	// }
}
