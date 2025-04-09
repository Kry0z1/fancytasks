package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Kry0z1/fancytasks/pkg/middleware"
	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to redis:", err.Error())
	}

	fmt.Println(pong)

	greet := func(w http.ResponseWriter, r *http.Request) error {
		return io.EOF
	}

	http.Handle("GET /", middleware.CollectErrorFunc(greet, middleware.Logger))

	fmt.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
