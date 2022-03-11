package app

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
	return StateManager{currentState, conf.Config.Subs, false, []Transform{}}
}

// Commit should be called to apply pending changes
func (s *StateManager) Commit() {
	if s.dirty {
		return
	}

	for _, transform := range s.pending {
		s.next = transform(s.next)
	}
	if currentState != s.next {
		currentState = s.next
	}
	Notify(s.next)
	s.dirty = true
}

// QueueTransform allows changes of state to be expressed and supplied as functions
// that take a state value and return a transformed value
func (s *StateManager) QueueTransform(transform Transform) {
	if s.dirty {
		return
	}
	s.pending = append(s.pending, transform)
}

// Notify will prompt all controllers to check if they need to update and if so, they will queue
// an update
func Notify(s State) {
	for _, receiver := range receivers {
		receiver.Receive(s)
	}
}
