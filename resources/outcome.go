package resources

// State represents one of two possible service states: UP (if healthy) and DOWN (if unhealthy)
type State string

const (
	StateUp   State = "UP"   // StateUp represents a healthy state of a service
	StateDown State = "DOWN" // StateDown represents an unhealthy state of a service
)

// BoolToState converts boolean type to a State type: if b is true, State is StateUp, otherwise State is StateDown
func BoolToState(b bool) State {
	return map[bool]State{
		true:  StateUp,
		false: StateDown,
	}[b]
}
