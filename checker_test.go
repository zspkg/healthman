package healthman

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHealthyChecker_Get(t *testing.T) {
	checker := NewHealthyChecker(
		&HealthmanSettings{Period: time.Second},
		nil,
		NewTestServiceHealthy(),
		NewTestServiceUnhealthy())
	checker.Run(context.Background())
	health, err := checker.Get(testServiceHealthyName)
	assert.Nil(t, err)
	assert.Equal(t, *health, ServiceHealth{
		ServiceName: testServiceHealthyName,
		Healthy:     true,
	})

	health, err = checker.Get(testServiceUnhealthyName)
	assert.Nil(t, err)
	assert.Equal(t, *health, ServiceHealth{
		ServiceName: testServiceUnhealthyName,
		Healthy:     false,
	})
}

func TestHealthyChecker_Info(t *testing.T) {
	var (
		healthyService   = NewTestServiceHealthy()
		unhealthyService = NewTestServiceUnhealthy()
	)
	checker := NewHealthyChecker(
		&HealthmanSettings{Period: time.Millisecond * 100},
		nil,
		healthyService,
		unhealthyService)
	checker.Run(context.Background())

	healthyService.ChangeHealthStatus()
	unhealthyService.ChangeHealthStatus()

	time.Sleep(time.Second)
	health := checker.Info()
	assert.Equal(t, health, []ServiceHealth{
		{
			ServiceName: testServiceHealthyName,
			Healthy:     false,
		},
		{
			ServiceName: testServiceUnhealthyName,
			Healthy:     true,
		},
	})
}
