package healthman

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/kit/kv"
	"testing"
	"time"
)

const (
	testConfigName      = ".config.test.yaml"
	testConfigWrongName = ".some.random.yaml"
)

func TestHealthy_Get(t *testing.T) {
	t.Run("correct config name: must succeed", func(t *testing.T) {
		configer := NewHealthmanConfiger(kv.NewViperFile(testConfigName))
		cfg := configer.HealthmanConfig()
		assert.Equal(t, cfg.Period, 2*time.Second)
	})

	t.Run("wrong config name: must fail", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("This code did not panic")
			}
		}()

		configer := NewHealthmanConfiger(kv.NewViperFile(testConfigWrongName))
		_ = configer.HealthmanConfig()
	})
}
