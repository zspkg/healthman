package healthman

import (
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisHealthier_Name(t *testing.T) {
	redis, _ := redismock.NewClientMock()
	healthier := NewRedisHealthier(redis)
	assert.Equal(t, RedisDependencyName, healthier.Name())
}

func TestRedisHealthier_GetHealthy(t *testing.T) {
	t.Run("no error on ping", func(t *testing.T) {
		redis, mock := redismock.NewClientMock()
		mock.ExpectPing().SetVal("everything is ok")
		mock.ExpectPing().SetErr(nil)
		var (
			healthier = NewRedisHealthier(redis)
			healthy   = healthier.CheckHealth()
		)
		assert.Equal(t, true, healthy)
	})
	t.Run("error on ping", func(t *testing.T) {
		redis, mock := redismock.NewClientMock()
		mock.ExpectPing().SetErr(errors.New("some error"))

		var (
			healthier = NewRedisHealthier(redis)
			healthy   = healthier.CheckHealth()
		)
		assert.Equal(t, false, healthy)
	})
}
