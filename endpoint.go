package healthman

import (
	"context"
	"github.com/zspkg/healthman/resources"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

const healthCheckerCtxKey = "health_checker"

// HealthCheckerSetter is a function returning a function that adds a HealthChecker to a context
func HealthCheckerSetter(checker HealthyChecker) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, healthCheckerCtxKey, checker)
	}
}

// HealthCheck returns full information about services' health based on HealthyChecker
// specified in context (by calling an Info function)
func HealthCheck(r *http.Request) []ServiceHealth {
	return r.Context().Value(healthCheckerCtxKey).(HealthyChecker).Info()
}

// CheckHealth returns a status OK and renders healthy status of all dependencies.
// To use this endpoint, one MUST include HealthCheckerSetter in context when initializing a router
func CheckHealth(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, formHealthResponse(HealthCheck(r)))
}

func formHealthResponse(servicesHealth []ServiceHealth) resources.ServiceHealthResponse {
	result := resources.ServiceHealthResponse{
		Outcome: resources.StateUp,
		Checks:  make([]resources.Check, len(servicesHealth)),
	}
	for i, serviceHealth := range servicesHealth {
		result.Checks[i] = serviceHealthToResource(serviceHealth)
	}

	return result
}

func serviceHealthToResource(serviceHealth ServiceHealth) resources.Check {
	return resources.Check{
		State: resources.BoolToState(serviceHealth.Healthy),
		Name:  serviceHealth.ServiceName,
	}

}
