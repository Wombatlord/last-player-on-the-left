package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/wombatlord/last-player-on-the-left/src/RSS"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
	"gopkg.in/yaml.v2"
)

var args struct {
	Alias string `arg:"positional" help:"The RSS feed alias"`
	Subscribe string `arg:"-s, --subscribe" help:"Supply a URL to the feed to create a subscription with the provided alias"`
}

type Subscription struct {
	Url string 
	Alias string
}

var (
	// url  = os.Args[1]
	feed RSS.RSSFeed
)

func main() {
	arg.MustParse(&args)
	sub, err := yaml.Marshal(
		Subscription{
			Url: args.Subscribe, 
			Alias: args.Alias,
		},
	)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	os.WriteFile("config.yaml", sub, 0644)

	content, err := RSS.GetContent(args.Subscribe)
	
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	episodeX := content.Channel[0].Item[1].Enclosure.Url
	fmt.Println(episodeX)

	feed.EpisodeData(*content)
	feed.EpisodeLink(*content, 10)
	lastplayer.StreamAudio("stream", episodeX)
}
