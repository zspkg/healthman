package healthman

import (
	"context"
	"errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"sync"
)

var conversionErr = errors.New("failed to convert to StoredHealthable type")

const healthCheckerProcessorName = "health_checker"

// StoredHealthable is a struct contained in InmemoryHealthyChecker storage
type StoredHealthable struct {
	Service Healthable
	Status  ServiceHealth
}

// healthyChecker is a struct that checks healthy status of a specified set of Healthable services
type healthyChecker struct {
	storage  *sync.Map
	settings HealthmanSettings
	log      *logan.Entry
}

// NewHealthyChecker creates a new instance of healthyChecker that implements HealthyChecker interface
func NewHealthyChecker(settings *HealthmanSettings, log *logan.Entry, services ...Healthable) HealthyChecker {
	storage := &sync.Map{}
	for _, service := range services {
		storage.Store(service.Name(), StoredHealthable{
			Service: service,
			Status: ServiceHealth{
				ServiceName: service.Name(),
				Healthy:     service.CheckHealth(),
			},
		})
	}
	if settings == nil {
		settings = &defaultHealthmanSettings
	}
	return &healthyChecker{
		storage:  storage,
		log:      log,
		settings: *settings,
	}
}

// Run launches a processor that every specified interval of time updates healthy status of services
func (c *healthyChecker) Run(ctx context.Context) {
	if c.log != nil {
		c.log.Infof("started %s", healthCheckerProcessorName)
	}
	go running.WithBackOff(
		ctx,
		c.log,
		healthCheckerProcessorName,
		c.Iterate,
		c.settings.Period,
		c.settings.Period,
		c.settings.Period,
	)
}

// Iterate ranges through all stored services in map and updates their healthy statuses
func (c *healthyChecker) Iterate(_ context.Context) error {
	c.storage.Range(func(key, value any) bool {
		stored := value.(StoredHealthable)
		c.storage.Store(key, StoredHealthable{
			Service: stored.Service,
			Status: ServiceHealth{
				ServiceName: stored.Service.Name(),
				Healthy:     stored.Service.CheckHealth(),
			},
		})
		return true
	})

	return nil
}

func (c *healthyChecker) Get(serviceName string) (*ServiceHealth, error) {
	v, ok := c.storage.Load(serviceName)
	if !ok {
		return nil, nil
	}

	stored, ok := v.(StoredHealthable)
	if !ok {
		return nil, conversionErr
	}

	return &stored.Status, nil
}

func (c *healthyChecker) Info() []ServiceHealth {
	result := make([]ServiceHealth, 0)
	c.storage.Range(func(_, value any) bool {
		stored := value.(StoredHealthable)
		result = append(result, ServiceHealth{
			ServiceName: stored.Status.ServiceName,
			Healthy:     stored.Status.Healthy,
		})
		return true
	})

	return result
}
