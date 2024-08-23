package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotExist = errors.New("key doesn't exist")
)

type Cache interface {
	// Set sets v to the key for duration of time.
	Set(key string, v string, duration time.Duration) error

	// Gets the value of key. returns ErrNotExist if key
	// Doesn't exist.
	Get(key string) (string, error)
}

type Redis struct {
	DB *redis.Client
}

func NewRedis() Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &Redis{
		DB: rdb,
	}
}

func (r *Redis) Set(key string, v string, duration time.Duration) error {
	err := r.DB.Set(context.Background(), key, v, duration).Err()
	return err
}

func (r *Redis) Get(key string) (string, error) {
	val, err := r.DB.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", ErrNotExist
	} else if err != nil {
		return "", err
	}

	return val, nil
}
