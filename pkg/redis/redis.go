package redis

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var rc *redis.Client

func init() {
	rc = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		log.Fatal("Failed to connect to redis:", err.Error())
	}
}
