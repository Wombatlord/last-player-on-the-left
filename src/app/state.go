package app

import (
	"fmt"
)

type State struct {
	FeedIndex    int
	EpisodeIndex int
	Initialised  bool
}

var currentState State

type Transform func(state State) State

type StateManager struct {
	next    State
	subs    []Subscription
	dirty   bool
	pending []Transform
	logger  chan string
}

var receivers []Receiver

type Receiver interface {
	Receive(s State)
}

// Register is used to connect a controller up for state update notifications
func Register(controller Receiver) Receiver {
	receivers = append(receivers, controller)
	return controller
}

// NewManager should be called to retrieve a manager instance, once commit has been called,
func NewManager() StateManager {
	if !currentState.Initialised {
		currentState.FeedIndex = 0
		currentState.EpisodeIndex = 0
		currentState.Initialised = true
	}
	return StateManager{
		currentState,
		conf.Config.Subs,
		false,
		[]Transform{},
		GetLogChan("StateManager"),
	}
}

// Commit should be called to apply pending changes
func (s *StateManager) Commit() {
	if !s.dirty {
		return
	}

	s.logger <- fmt.Sprintf("Applying %d transforms", len(s.pending))
	for _, transform := range s.pending {
		s.next = transform(s.next)
	}
	s.pending = []Transform{}
	if currentState != s.next {
		currentState = s.next
		s.Notify(s.next)
	}
	s.dirty = false
}

// QueueTransform allows changes of state to be expressed and supplied as functions
// that take a state value and return a transformed value
func (s *StateManager) QueueTransform(transform Transform) {
	s.pending = append(s.pending, transform)
	if !s.dirty {
		s.dirty = true
	}
}

// Notify will prompt all controllers to check if they need to update and if so, they will queue
// an update
func (s *StateManager) Notify(state State) {
	s.logger <- fmt.Sprintf("Notifying receivers of state: %+v", state)
	for _, receiver := range receivers {
		receiver.Receive(state)
	}
}
