package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"github.com/wombatlord/last-player-on-the-left/src/lastplayer"
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
	feed   *clients.RSSFeed
	logger chan string
	conf   *app.ConfigFile
	err    error
)

var AppControllers view.Controllers

// ConfigureUI Sets up the actual content
func ConfigureUI() *tview.Application {
	gui := tview.NewApplication()

	AppControllers.FeedMenu = view.NewFeedsController(gui)
	domain.Register(AppControllers.FeedMenu)

	AppControllers.EpisodeMenu = view.NewEpisodeMenuController()
	domain.Register(AppControllers.EpisodeMenu)

	AppControllers.APViewController = view.NewAPViewController()
	domain.Register(AppControllers.APViewController)

	AppControllers.RootController = view.NewRootController(gui)

	audioPanel := lastplayer.FetchAudioPanel()
	audioPanel.SubscribeToState(AppControllers.APViewController)
	audioPanel.AttachApp(gui)
	audioPanel.SpawnPublisher()

	application := view.Build(gui, AppControllers)

	return application
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
		log.Fatal(ConfigureUI().Run())
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
}

func fatal(err error) {
	if err != nil {
		logger <- fmt.Sprintf("error: %v", err)
		log.Fatalf("error: %v", err)
	}
}
