package cache

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sinulingga23/form-builder-be/define"
)

var (
	client *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Network:  os.Getenv("REDIS_SERVER_NETWORK"),
		Password: os.Getenv("REDIS_SERVER_PASSWORD"),
		DB:       0,
	})
}

func GetValue(ctx context.Context, key string) (interface{}, error) {
	stringCmd := client.Get(ctx, key)
	return stringCmd, stringCmd.Err()
}

func SetValue(ctx context.Context, key string, value interface{}, duration ...time.Duration) error {
	if len(duration) > 0 {
		return client.Set(ctx, key, value, duration[0]).Err()
	}

	return client.Set(ctx, key, value, define.DEFAULT_TIME_TO_LIVE_CACHE*time.Hour).Err()
}

func DelValue(ctx context.Context, key ...string) error {
	return client.Del(ctx, key...).Err()
}
