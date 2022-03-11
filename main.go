package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
	"github.com/wombatlord/last-player-on-the-left/src/view"
	"log"
	"os"
)

var args struct {
	Alias        string `arg:"positional" help:"The RSS feed alias"`
	Subscription string `arg:"-s, --subscribe" help:"Supply a URL to the feed to create a subscription with the provided alias"`
	Latest       bool   `arg:"-l, --latest" help:"Play the latest episode associated to the alias"`
	Episode      int    `arg:"-e, --episode" help:"Play a specific episode. 0 is the latest episode."`
}

var (
	// url  = os.Args[1]
	feed   *clients.RSSFeed
	logger chan string
	conf   *app.ConfigFile
	err    error
)

func playAudio(url string) {
	lastplayer.StreamAudio("stream", url)
}

var AppControllers view.Controllers

// ConfigureUI Sets up the actual content
func ConfigureUI() *tview.Application {
	AppControllers.FeedMenu = view.NewFeedsController()
	app.Register(AppControllers.FeedMenu)

	AppControllers.EpisodeMenu = view.NewEpisodeMenuController()
	app.Register(AppControllers.EpisodeMenu)

	return view.Build(AppControllers)
}

func mainLogger() chan string {
	return app.GetLogChan("main")
}

func main() {
	// Create the logger
	logger = app.GetLogChan("main")
	defer close(logger)

	// Parse the args
	arg.MustParse(&args)
	logger <- fmt.Sprintf("Args parsed: %+v", args)

	// Load the config file
	conf, err = app.LoadConfig("config.yaml")
	fatal(err)

	if len(os.Args) == 1 {
		gui := ConfigureUI()
		log.Fatal(gui.Run())
	}

	// Pull the feed
	if args.Subscription != "" {
		feed, err = clients.GetContent(args.Subscription)
	} else {
		sub := conf.Config.GetByAlias(args.Alias)
		if sub.Url == "" {
			fmt.Printf("You have no subscription for the alias %s\n", args.Alias)
		}
		feed, err = clients.GetContent(sub.Url)
	}
	fatal(err)

	// If no episode arg is provided, set value to -1 to prevent 0 value instantiation.
	args.Episode = -1
	// Prints data from the feed for sanity checking.
	// feed.EpisodeData(*feed)

	// Main Control flow
	if args.Subscription != "" {
		// Handles new Subs
		fatal(conf.Include(args.Alias, args.Subscription))
	}

	if args.Latest {
		// Get the url for the latest episode in the feed.
		latest := feed.Channel[0].Item[0]
		streamUrl := latest.Enclosure.Url

		// do the latest noises
		playAudio(streamUrl)
	}

	if args.Episode > -1 {
		// Get the url for the requested episode in the feed.
		episode := feed.Channel[0].Item[args.Episode]
		streamUrl := episode.Enclosure.Url

		// do the specific noises
		playAudio(streamUrl)
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
