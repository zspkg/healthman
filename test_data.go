package healthman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"io"
	"net/http"
	"sync"
	"testing"
)

const (
	testServiceHealthyName   = "service_healthy"
	testServiceUnhealthyName = "service_unhealthy"
)

type TestService struct {
	healthy bool
	name    string
	mutex   sync.Mutex
}

func NewTestServiceHealthy() *TestService {
	return &TestService{healthy: true, name: testServiceHealthyName, mutex: sync.Mutex{}}
}

func NewTestServiceUnhealthy() *TestService {
	return &TestService{healthy: false, name: testServiceUnhealthyName, mutex: sync.Mutex{}}
}

func (t *TestService) Name() string {
	t.mutex.Lock()
	result := t.name
	t.mutex.Unlock()
	return result
}

func (t *TestService) CheckHealth() bool {
	t.mutex.Lock()
	result := t.healthy
	t.mutex.Unlock()
	return result
}

func (t *TestService) ChangeHealthStatus() {
	t.mutex.Lock()
	t.healthy = !t.healthy
	t.mutex.Unlock()
}

type testCli struct {
	t     *testing.T
	token string
}

func NewTestClient(authToken string, t *testing.T) *testCli {
	return &testCli{
		token: authToken,
		t:     t,
	}
}

func (c *testCli) Do(method, endpoint string, data []byte) *http.Response {
	// creating request
	bodyReader := bytes.NewReader(data)
	request, err := http.NewRequest(method, endpoint, bodyReader)
	if err != nil {
		c.t.Fatal(errors.Wrap(err, "failed to create connector request"))
	}

	// setting headers
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// sending request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		c.t.Fatal(errors.Wrap(err, "failed to process request"))
	}
	if response == nil {
		c.t.Fatal(errors.New("failed to process request: response is nil"))
	}

	return response
}

func (c *testCli) ReadBody(response *http.Response, dst interface{}) {
	defer func(body io.Closer) {
		err := body.Close()
		if err != nil {
			c.t.Fatal(err)
		}
	}(response.Body)

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		c.t.Fatal(errors.Wrap(err, "failed to read response body"))
	}

	err = json.Unmarshal(raw, &dst)
	if err != nil {
		c.t.Fatal(err)
	}
}
