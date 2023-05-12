package healthman

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/zspkg/healthman/resources"
	"gitlab.com/distributed_lab/ape"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheckEndpoint(t *testing.T) {
	var (
		endpoint = "/integrations/some-service/health"
		client   = NewTestClient("some_token", t)
		router   = chi.NewRouter()
		checker  = NewHealthyChecker(nil, nil, NewTestServiceHealthy(), NewTestServiceUnhealthy())
	)
	router.Use(ape.CtxMiddleware(HealthCheckerSetter(checker)))
	router.Get(endpoint, CheckHealth)

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	time.Sleep(time.Second)

	t.Run("must get healthy and unhealthy services", func(t *testing.T) {
		t.Log(testServer.URL + endpoint)
		response := client.Do(http.MethodGet, testServer.URL+endpoint, nil)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		var healthResponse resources.ServiceHealthResponse
		client.ReadBody(response, &healthResponse)
		assert.Equal(t, true, Equal([]resources.Check{
			{
				State: resources.StateUp,
				Name:  testServiceHealthyName,
			},
			{
				State: resources.StateDown,
				Name:  testServiceUnhealthyName,
			},
		}, healthResponse.Checks))
	})
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal[K comparable](a, b []K) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		found := false
		for j := range b {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
