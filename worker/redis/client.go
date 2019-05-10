package redis

import (
	"saas/env"

	"github.com/go-redis/redis"
)

// Client pointer for redis
func Client() *redis.Client {
	options := &redis.Options{
		Addr: env.Getenv("REDIS_HOST", "localhost:6379"),
		// TODO: Use secrets in docker-compose
		Password: env.Getenv("REDIS_PASSWORD", ""), // no password set
		DB:       0,                                // use default DB
	}
	return redis.NewClient(options)
}
