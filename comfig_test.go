package healthman

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/kit/kv"
	"testing"
	"time"
)

const testConfigName = ".config.test.yaml"

func TestHealthy_Get(t *testing.T) {
	configer := NewHealthmanConfiger(kv.NewViperFile(testConfigName))
	cfg := configer.HealthmanConfig()
	assert.Equal(t, cfg.Period, 2*time.Second)
}
