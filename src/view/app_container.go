package view

import (
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/audiopanel"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"log"
	"os"
)

// Views is the declaration of the full set of views that must be supplied
// to the LastPlayer on Build
type Views struct {
	Root        *tview.Flex
	TopRow      *tview.Flex
	EpisodeMenu *tview.List
	FeedMenu    *tview.List
	APView      *tview.TextView
}

// Controllers is the declaration of the full set of controllers
// that must be supplied to the app on Build
type Controllers struct {
	FeedMenu         *FeedsMenuController
	EpisodeMenu      *EpisodeMenuController
	RootController   *RootController
	APViewController *APViewController
}

// LastPlayer extends the tview.Application with our custom functionality
type LastPlayer struct {
	*tview.Application
	Controllers Controllers
	Views       Views
	FocusRing   []tview.Primitive
	State       *domain.State
	AudioPanel  *audiopanel.AudioPanel
	Config      app.Config
	LogFile     *os.File
}

// subscribePanelAware is used to subscribe PanelStateAwareController instances to
// periodic player state updates
func (lp *LastPlayer) subscribePanelAware(subs ...PanelStateAwareController) {
	for _, aware := range subs {
		lp.AudioPanel.SubscribeToState(aware)
	}
}

// registerReceivers is used to register ReceiverController instances to be notified
// of state changes
func (lp *LastPlayer) registerReceivers(receivers ...ReceiverController) {
	for _, receiver := range receivers {
		domain.Register(receiver)
	}
}

// declareFocusRing is used to indicate to LastPlayer which Primitives should be cycled through
// on tab key press
func (lp *LastPlayer) declareFocusRing(views ...tview.Primitive) {
	lp.FocusRing = views
}

// Build and returns the LastPlayer, implement any additions to
// the user interface adding new primitives to the lastPlayer hierarchy
// in this function
func Build() *LastPlayer {
	config, _ := app.LoadConfig()
	application := &LastPlayer{
		Application: tview.NewApplication(),
		Config:      config.Config,
		State:       (&domain.State{}).Init(),
	}
	application.AudioPanel = audiopanel.
		FetchAudioPanel().
		AttachApp(application)

	logfile, _ := os.Open(application.Config.Logs)
	application.LogFile = logfile
	log.SetOutput(logfile)

	application.Views = Views{
		Root:        MainFlex(),
		TopRow:      TopRow(),
		EpisodeMenu: EpisodeMenu(),
		FeedMenu:    FeedMenu(),
		APView:      AudioPanelView(),
	}
	application.SetRoot(application.Views.Root, true)

	application.Controllers = Controllers{
		FeedMenu:         NewFeedsController(application),
		EpisodeMenu:      NewEpisodeMenuController(application),
		APViewController: NewAPViewController(application),
		RootController:   NewRootController(application),
	}

	application.registerReceivers(
		application.Controllers.EpisodeMenu,
		application.Controllers.APViewController,
	)

	application.subscribePanelAware(
		application.Controllers.APViewController,
	)

	application.declareFocusRing(
		application.Views.FeedMenu,
		application.Views.EpisodeMenu,
	)

	application.setupLayout()

	return application
}

// Run overrides the tview.Application Run method and includes a deferred close
// of the logfile
func (lp *LastPlayer) Run() error {
	defer func(LogFile *os.File) {
		err := LogFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(lp.LogFile)
	return lp.Application.Run()
}

// GetLogger can be used to get a log.Logger with the prefix as passed. This
// can be accessed inside controllers etc.
func (lp *LastPlayer) GetLogger(prefix string) *log.Logger {
	return log.New(lp.LogFile, prefix, 0)
}

// setupLayout manages the nesting and sizes of the various views
func (lp *LastPlayer) setupLayout() {
	lp.Views.TopRow.AddItem(lp.Views.FeedMenu, -1, 1, true)
	lp.Views.TopRow.AddItem(lp.Views.EpisodeMenu, -1, 1, true)

	lp.Views.Root.AddItem(lp.Views.TopRow, -1, 4, true)
	lp.Views.Root.AddItem(lp.Views.APView, -1, 1, false)
}

// QueueUpdate overrides the tview.Application method QueueUpdate by
// wrapping the call to the passed function in a closure that checks if
// it updated the global state and if so, calls domain.Notify with the
// updated state value
func (lp *LastPlayer) QueueUpdate(f func()) *LastPlayer {
	closure := func() {
		oldState := *lp.State
		f()
		if oldState != *lp.State {
			domain.Notify(*lp.State)
		}
	}
	lp.Application.QueueUpdate(closure)
	return lp
}
