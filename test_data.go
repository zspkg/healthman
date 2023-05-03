package healthman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"io"
	"net/http"
	"testing"
)

const (
	testServiceHealthyName   = "service_healthy"
	testServiceUnhealthyName = "service_unhealthy"
)

type testServiceHealthy struct{}
type testServiceUnhealthy struct{}

func NewTestServiceHealthy() Healthable {
	return testServiceHealthy{}
}

func (t testServiceHealthy) Name() string {
	return testServiceHealthyName
}

func (t testServiceHealthy) CheckHealth() bool {
	return true
}

func NewTestServiceUnhealthy() Healthable {
	return testServiceUnhealthy{}
}

func (t testServiceUnhealthy) Name() string {
	return testServiceUnhealthyName
}

func (t testServiceUnhealthy) CheckHealth() bool {
	return false
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
