package healthman

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const RedisDependencyName = "redis"

type redisHealthier struct {
	client *redis.Client
}

// NewRedisHealthier is a constructor of a redis database healthier that implements CheckHealth method based
// on provided redis connection
func NewRedisHealthier(client *redis.Client) Healthable {
	return &redisHealthier{
		client: client,
	}
}

// CheckHealth pings Redis and returns whether no error occurred
func (h *redisHealthier) CheckHealth() bool {
	_, err := h.client.Ping(context.Background()).Result()
	return err == nil
}

// Name returns a determined RedisDependencyName value
func (h *redisHealthier) Name() string {
	return RedisDependencyName
}
