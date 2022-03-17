package main

import (
	"github.com/alexflint/go-arg"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/view"
	"log"
	"os"
)

var args struct {
	Alias        string `arg:"positional" help:"The RSS feed alias"`
	Subscription string `arg:"-s, --subscribe" help:"Supply a URL to the feed to create a subscription with the provided alias"`
	Download     []int  `arg:"-d, --download" help:"Download episodes. 0 is the latest episode."`
}

var (
	conf *app.ConfigFile
	err  error
)

func main() {
	// Load the config file
	conf, err = app.LoadConfig()
	fatal(err)

	// Create the logger
	logfile, err := os.Open(conf.Config.Logs)
	fatal(err)
	logger := log.New(logfile, "main", 0)

	// Parse the args
	arg.MustParse(&args)
	logger.Printf("Args parsed: %+v", args)

	// If invoked bare, start the main app
	if len(os.Args) == 1 {
		lastPlayer := view.Build()
		fatal(lastPlayer.Run())
	}

	// Add subscription if -s flag supplied
	if args.Subscription != "" {
		fatal(conf.Include(args.Alias, args.Subscription))
		fatal(conf.Save())
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
