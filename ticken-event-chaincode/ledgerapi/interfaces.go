package ledgerapi

// StateList functions that a state list
// should have
type StateList interface {
	// AddState puts state into world state
	AddState(state State) error

	// UpdateState puts state into world state. Same as AddState but
	// separate as semantically different
	UpdateState(state State) error

	// GetState returns state from world state. Unmarshalls the JSON
	// into passed state. Key is the split key value used in Add/Update
	// joined using a col
	GetState(key string, state State) error

	// StateExists returns true if the element exists in the world state.
	// Error should be nil if the element don't exist.
	StateExists(key string) (bool, error)
}

// State interface states must implement
// for use in a list
type State interface {
	GetKey() string
	Serialize() ([]byte, error)
}
