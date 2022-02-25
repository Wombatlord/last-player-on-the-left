package main

import (
	"fmt"

	"github.com/wombatlord/last-player-on-the-left/src/app"

	"github.com/alexflint/go-arg"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
	"github.com/wombatlord/last-player-on-the-left/src/rss"
)

var args struct {
	Alias     string `arg:"positional" help:"The RSS feed alias"`
	Subscription string `arg:"-s, --subscribe" help:"Supply a URL to the feed to create a subscription with the provided alias"`
	Latest bool `arg:"-l, --latest" help:"Play the latest episode associated to the alias"`
}

var (
	// url  = os.Args[1]
	feed *RSS.RSSFeed
)

func playAudio(url string) {
	lastplayer.StreamAudio("stream", url)
}

func main() {
	// Load the config file
	conf, err := app.LoadConfig("config.yaml")
	fatal(err)
	
	// Parse the args
	arg.MustParse(&args)

	// Pull the feed
	if args.Subscription != "" {
		feed, err = RSS.GetContent(args.Subscription)
	} else {
		feed, err = RSS.GetContent(conf.Subs[args.Alias])
	}
	fatal(err)
	
	// Prints data from the feed for sanity checking.
	// feed.EpisodeData(*feed)
	
	// Main Control flow
	if args.Subscription != "" {
		// Handles new Subs
		fatal(conf.Include(args.Alias, args.Subscription))
	}
	
	if args.Latest {
		// Play the latest episode in the feed.
		latest := feed.Channel[0].Item[0]
		streamUrl := latest.Enclosure.Url
		fmt.Printf("URL ===> %s\n\n\n\nTITLE 8===D %s", streamUrl, latest.Title)
		
		// do the noises
		playAudio(streamUrl)
	}
}

func fatal(err error) {
	if err != nil {
		panic(err)
		// log.Fatalf("error: %v", err)
	}
}