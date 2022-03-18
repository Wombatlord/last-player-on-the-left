package domain

import (
	"github.com/wombatlord/last-player-on-the-left/src/clients"
)

// State represents all the shared global application state that is not managed by the
// audiopanel
type State struct {
	FeedIndex      int
	Feed           *clients.RSSFeed
	EpisodeIndex   int
	PlayingEpisode *clients.Item
	Initialised    bool
}

// Init initialises the state
func (s *State) Init() *State {
	s.FeedIndex = NoItem
	s.Feed = nil
	s.PlayingEpisode = nil
	s.Initialised = true

	return s
}

// NoItem is used to indicate a "none" value for a menu index (i.e. a positive integer)
const NoItem int = -1

var receivers []Receiver

// Receiver is the interface that must be implemented to be eligible for state change
// notification
type Receiver interface {
	Receive(s State)
}

// Register is used to connect a controller up for state update notifications
func Register(controller Receiver) Receiver {
	receivers = append(receivers, controller)
	return controller
}

// Notify will prompt all controllers to check if they need to update and if so, they will queue
// an update
func Notify(state State) {
	for _, receiver := range receivers {
		receiver.Receive(state)
	}
}
