package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wombatlord/last-player-on-the-left/src/app"
	"github.com/wombatlord/last-player-on-the-left/src/audiopanel"
	"github.com/wombatlord/last-player-on-the-left/src/clients"
	"github.com/wombatlord/last-player-on-the-left/src/domain"
	"io/fs"
	"log"
	"os"
)

type BeforeDraw func(_ tcell.Screen) bool

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
	Controllers   Controllers
	Views         Views
	FocusRing     []tview.Primitive
	State         *domain.State
	AudioPanel    *audiopanel.AudioPanel
	Config        app.Config
	LogFile       *os.File
	logger        *log.Logger
	previousState domain.State
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
	initialState := (&domain.State{}).Init()
	application := &LastPlayer{
		Application:   tview.NewApplication(),
		Config:        config.Config,
		State:         initialState,
		previousState: *initialState,
	}
	application.AudioPanel = audiopanel.
		FetchAudioPanel().
		AttachLogger(application.GetLogger("AudioPanel"))
	application.AudioPanel.
		SetPublishCallback(
			func(f func()) { application.QueueUpdateDraw(f) },
		)

	logfile, _ := os.OpenFile(application.Config.Logs, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	application.LogFile = logfile
	log.SetOutput(logfile)
	clients.InitLoggers(application.GetLogger)

	application.Views = Views{
		Root:        MainFlex(),
		TopRow:      TopRow(),
		EpisodeMenu: EpisodeMenu(),
		FeedMenu:    FeedMenu(),
		APView:      AudioPanelView(),
	}

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
	application.SetRoot(application.Views.Root, true)
	application.SetBeforeDrawFunc(application.notifyCheck())

	return application
}

// notifyCheck compares the state before the last draw and if it
// has changed, notify is called
func (lp *LastPlayer) notifyCheck() BeforeDraw {
	return func(_ tcell.Screen) bool {
		s := *lp.State
		if lp.previousState != s {
			domain.Notify(s)
		}
		lp.previousState = s

		return false
	}
}

// Run overrides the tview.Application Run method and includes a deferred close
// of the logfile
func (lp *LastPlayer) Run() error {
	_ = os.MkdirAll(lp.Config.Cache, fs.ModeDir+fs.FileMode(0774))
	lp.AudioPanel.SpawnPublisher()
	defer func(LogFile *os.File) {
		err := LogFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		_ = os.RemoveAll(lp.Config.Cache)
		_ = os.MkdirAll(lp.Config.Cache, fs.ModeDir+fs.FileMode(0774))
	}(lp.LogFile)
	return lp.Application.Run()
}

// GetLogger can be used to get a log.Logger with the prefix as passed. This
// can be accessed inside controllers etc.
func (lp *LastPlayer) GetLogger(prefix string) *log.Logger {
	return log.New(lp.LogFile, "[ "+prefix+" ]: ", 0)
}

// setupLayout manages the nesting and sizes of the various views
func (lp *LastPlayer) setupLayout() {
	lp.Views.TopRow.AddItem(lp.Views.FeedMenu, -1, 1, true)
	lp.Views.TopRow.AddItem(lp.Views.EpisodeMenu, -1, 1, true)

	lp.Views.Root.AddItem(lp.Views.TopRow, -1, 4, true)
	lp.Views.Root.AddItem(lp.Views.APView, -1, 1, false)
}
