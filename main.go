package main

import (
	"fmt"

	"github.com/wombatlord/last-player-on-the-left/src/app"

	"github.com/alexflint/go-arg"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
)

var args struct {
	Alias     string `arg:"positional" help:"The RSS feed alias"`
	Subscription string `arg:"-s, --subscribe" help:"Supply a URL to the feed to create a subscription with the provided alias"`
	Latest bool `arg:"-l, --latest" help:"Play the latest episode associated to the alias"`
	Episode int `arg:"-e, --episode" help:"Play a specific episode. 0 is the latest episode."`
}

var (
	// url  = os.Args[1]
	feed *clients.RSSFeed
)

func playAudio(url string) {
	lastplayer.StreamAudio("stream", url)
}

func playLatest() {
	// Get the url for the latest episode in the feed.
	latest := feed.Channel[0].Item[0]
	streamUrl := latest.Enclosure.Url
	
	// do the latest noises
	playAudio(streamUrl)
}

func main() {
	// Load the config file
	conf, err := app.LoadConfig("config.yaml")
	fatal(err)
	
	// Parse the args
	arg.MustParse(&args)

	// Pull the feed
	if args.Subscription != "" {
		feed, err = clients.GetContent(args.Subscription)
	} else {
		feed, err = clients.GetContent(conf.Subs[args.Alias])
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
		playLatest()
	}

	if args.Episode != 0 {
		// Get the url for the requested episode in the feed.
		episode := feed.Channel[0].Item[args.Episode]
		streamUrl := episode.Enclosure.Url

		// do the specific noises
		playAudio(streamUrl)
	} else {
		fmt.Println("Please provide an episode number 1 or greater.\n 0 is the latest episode, use -l")
	}
}

func fatal(err error) {
	if err != nil {
		panic(err)
		// log.Fatalf("error: %v", err)
	}
}
