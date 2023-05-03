package healthman

import "context"

// ServiceHealth is a struct representing a health state of a service
type ServiceHealth struct {
	ServiceName string
	Healthy     bool
}

// Healthable is an interface that any service whose health can be checked must implement
type Healthable interface {
	Name() string
	CheckHealth() bool
}

// HealthyChecker is a processor interface that is responsible for pinging all services once in a while
// to check for their healthy statuses
type HealthyChecker interface {
	// Run is a function that starts a goroutine that checks once in a specified amount of time
	// health status of services
	Run(ctx context.Context)
	// Get is a function that, based on serviceName, returns a service health status
	Get(serviceName string) (*ServiceHealth, error)
	// Info is a function that returns full information about all services' health statuses
	Info() []ServiceHealth
}
