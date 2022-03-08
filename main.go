package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
	"log"
)

var args struct {
	Alias        string `arg:"positional" help:"The RSS feed alias"`
	Subscription string `arg:"-s, --subscribe" help:"Supply a URL to the feed to create a subscription with the provided alias"`
	Latest       bool   `arg:"-l, --latest" help:"Play the latest episode associated to the alias"`
	Episode      int    `arg:"-e, --episode" help:"Play a specific episode. 0 is the latest episode."`
	Download     []int  `arg:"-d, --download" help:"Download episodes. 0 is the latest episode."`
}

var (
	feed   *clients.RSSFeed
	logger chan string
	pics   []string
)

func playAudio(url string) {
	lastplayer.StreamAudio("stream", url)
}

func main() {
	// Create the logger
	logger = app.GetLogChan("main")
	defer close(logger)

	pics = []string{
		"https://decider.com/wp-content/uploads/2016/03/feature-henry-zebrowski-the-characters.jpg",
		"https://i0.wp.com/brightestyoungthings.com/wp-content/uploads/2018/08/Marcus.jpg",
		"https://d2eehagpk5cl65.cloudfront.net/img/c2400x1350-w2400-q80/uploads/2017/06/14968109224555.jpg",
	}

	downloader := clients.DownloadClient{Client: clients.NewClient()}
	reqs := downloader.CreateRequests(pics)

	downloader.DownloadMulti(*downloader.Client, reqs...)

	// for req := range reqs {
	// 	downloader.DownloadEpisode(*downloader.Client, reqs[req])
	// }

	// Parse the args
	arg.MustParse(&args)
	logger <- fmt.Sprintf("Args parsed: %+v", args)

	// Load the config file
	conf, err := app.LoadConfig("config.yaml")
	fatal(err)

	// If no episode arg is provided, set value to -1 to prevent 0 value instantiation.
	args.Episode = -1

	// Pull the feed
	if args.Subscription != "" {
		feed, err = clients.GetContent(args.Subscription)
	} else {
		url := conf.Config.Subs[args.Alias]
		if url == "" {
			fmt.Printf("You have no subscription for the alias %s\n", args.Alias)
		}
		feed, err = clients.GetContent(url)
	}
	fatal(err)

	// Prints data from the feed for sanity checking.
	// feed.EpisodeData(*feed)

	// Main Control flow
	if args.Subscription != "" {
		// Handles new Subs
		fatal(conf.Include(args.Alias, args.Subscription))
	}

	// Currently must be passed with a -l or -e to ensure player stays open while downloading.
	// If ints are passed with -d
	// Instantiate an array of strings: eps
	// Iterate over the ints and grab the urls to the associated episodes
	// Append each url to the eps array. End loop
	// Create a grab.Request for each episode url: reqs
	// Initiate a download for each request.
	if len(args.Download) > 0 {
		var eps []string
		logger <- "-d > 0"
		for ep := range args.Download {
			episode := feed.Channel[0].Item[args.Download[ep]]
			streamUrl := episode.Enclosure.Url
			eps = append(eps, streamUrl)
		}
		for ep := range eps {
			logger <- eps[ep]
		}
		downloads := downloader.CreateRequests(eps)
		downloader.DownloadMulti(*downloader.Client, downloads...)
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
		logger <- fmt.Sprintf("error: %v", err)
		log.Fatalf("error: %v", err)
	}
}
