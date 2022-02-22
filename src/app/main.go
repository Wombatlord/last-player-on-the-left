package main

import (
	"fmt"	
	"log"
	"os"

	"github.com/wombatlord/lastplayerontheleft/src/lastplayer"
	"github.com/wombatlord/lastplayerontheleft/src/rss"
)

var (
	url  = os.Args[1]
	feed RSS.RSSFeed
)

func main() {
	content, err := RSS.GetContent(url)
	
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	episodeX := content.Channel[0].Item[1].Enclosure.Url
	fmt.Println(episodeX)

	feed.EpisodeData(*content)
	feed.EpisodeLink(*content, 10)
	lastplayer.StreamAudio("stream", episodeX)
}
