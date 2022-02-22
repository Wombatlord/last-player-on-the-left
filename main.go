package main

import (
	"fmt"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
	"github.com/wombatlord/last-player-on-the-left/src/rss"
	"gopkg.in/yaml.v2"
)

var args struct {
	Alias     string `arg:"positional" help:"The RSS feed alias"`
	Subscribe string `arg:"-s, --subscribe" help:"Supply a URL to the feed to create a subscription with the provided alias"`
}

var (
	// url  = os.Args[1]
	feed RSS.RSSFeed
)

func main() {
	arg.MustParse(&args)
	if args.Subscribe != "" {
		sub, err := yaml.Marshal(
			app.Subscription{
				Url:   args.Subscribe,
				Alias: args.Alias,
			},
		)
		fatal(err)

		fatal(os.WriteFile("subscription.yaml", sub, 0644))
	}

	content, err := RSS.GetContent(args.Subscribe)
	fatal(err)

	episodeX := content.Channel[0].Item[1].Enclosure.Url
	fmt.Println(episodeX)

	feed.EpisodeData(*content)
	feed.EpisodeLink(*content, 10)
	lastplayer.StreamAudio("stream", episodeX)
}

func fatal(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
