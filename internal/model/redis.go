package model

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrNoRedisKey = errors.New("key not found in redis")

func newRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
