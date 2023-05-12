package healthman

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/zamslb/healthman/resources"
	"gitlab.com/distributed_lab/ape"
	"net/http"
	"net/http/httptest"
	"testing"
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

	t.Run("must get healthy and unhealthy services", func(t *testing.T) {
		t.Log(testServer.URL + endpoint)
		response := client.Do(http.MethodGet, testServer.URL+endpoint, nil)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		var healthResponse resources.ServiceHealthListResponse
		client.ReadBody(response, &healthResponse)
		assert.Equal(t, []resources.ServiceHealth{
			{
				Key: resources.Key{ID: "", Type: resources.SERVICE_HEALTH},
				Attributes: resources.ServiceHealthAttributes{
					Healthy:     true,
					ServiceName: "service_healthy",
				},
			},
			{
				Key: resources.Key{ID: "", Type: resources.SERVICE_HEALTH},
				Attributes: resources.ServiceHealthAttributes{
					Healthy:     false,
					ServiceName: "service_unhealthy",
				},
			},
		}, healthResponse.Data)
	})
}
