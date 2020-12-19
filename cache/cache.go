package cache

import (
	"carcompare/config"

	"github.com/go-redis/redis/v8"
)

// RedisDB redis connector
var RedisDB *redis.Client

// Connect to redis
func Connect() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     config.Config("REDIS_HOST") + ":" + config.Config("REDIS_PORT"),
		Password: config.Config("REDIS_PASSWORD"),
		DB:       0,
	})
}
