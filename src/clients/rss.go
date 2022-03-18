package clients

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type FeedCache map[string]*RSSFeed

var (
	feedCache FeedCache = make(map[string]*RSSFeed)
)

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

// Defining Structs to parse clients Feed xml from HTTP request.
// The full feed including header.
type RSSFeed struct {
	XMLName xml.Name  `xml:"rss"`
	Channel []Channel `xml:"channel"`
}

// Channel data represents an entire podcast.
type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Item        []Item   `xml:"item"`
	Generator   string   `xml:"generator"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Language    string   `xml:"language"`
	PubDate     string   `xml:"pubDate"`
}

// Item contains data for individual episodes.
type Item struct {
	XMLName     xml.Name  `xml:"item"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	PubDate     string    `xml:"pubDate"`
	Author      string    `xml:"author"`
	Link        string    `xml:"link"`
	Enclosure   Enclosure `xml:"enclosure"`
}

// GetContent retrieves a clients Feed via HTTP Request.
// Parse the xml in the response into structs.
// Exit & Print in the event of an error.
func GetContent(url string) (*RSSFeed, error) {
	loggers[RSSLog] = loggers[RSSLog]
	loggers[RSSLog].Printf("Retrieving RSS Feed at: %s", url)
	if feed, ok := feedCache[url]; !ok {
		loggers[RSSLog].Print("Cache Miss")
		feed = &RSSFeed{}
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("GET error: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				loggers[RSSLog].Fatal(err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status error: %v", resp.StatusCode)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %v", err)
		}

		err = xml.Unmarshal(data, feed)
		if err != nil {
			loggers[RSSLog].Printf("ERROR: %s", err.Error())
		}
		feedCache[url] = feed
	} else {
		loggers[RSSLog].Print("Cache Hit")
	}

	return feedCache[url], nil
}

// EpisodeData iterates over Item structs within the Channel struct.
// Print some episode information to the terminal.
func (RSSFeed *RSSFeed) EpisodeData(feed RSSFeed) {
	loggers[RSSLog].Print("Dumping episode data")
	for _, channel := range feed.Channel {
		for _, item := range channel.Item {
			fmt.Println("Title: " + item.Title)
			fmt.Println("Description: " + item.Description)
			fmt.Println("Author: " + item.Author)
			fmt.Println("Link: " + item.Link)
		}
	}
}

// EpisodeLink returns the link for a single episode.
func (RSSFeed *RSSFeed) EpisodeLink(feed RSSFeed, episodeNumber int) string {
	loggers[RSSLog].Printf("Retrieving link for episode number: %d", episodeNumber)
	episode := feed.Channel[0].Item[episodeNumber].Link
	return episode
}
