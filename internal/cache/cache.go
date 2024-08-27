package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotExist = errors.New("key doesn't exist")
	ErrBadValue = errors.New("value is not an int") // Used for Incr opertaion
)

type Cache interface {
	// Set sets v to the key for duration of time.
	Set(key string, v string, duration time.Duration) error

	// Get Get the value from cache using key
	//
	// Cache Hit returns a string and a nil.
	// Cache Miss returns "" and ErrNotExist.
	// Errors in Cache returns "", and
	Get(key string) (string, error)

	// Incr Increment value in the cache.
	//
	// if the key doesn't exist, Incr creates
	// the key and sets it to 1
	Incr(key string) (int, error)

	// Expire sets an expire time for a key
	// returns ErrNotExist in case of error
	Expire(key string, seconds int) error
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

func (r *Redis) Incr(key string) (int, error) {
	val, err := r.DB.Incr(context.Background(), key).Result()

	if err != nil {
		return 0, ErrBadValue
	}

	return int(val), nil
}

func (r *Redis) Expire(key string, seconds int) error {
	val, err := r.DB.Expire(context.Background(), key, time.Duration(seconds)*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrNotExist
		}
	}

	if !val {
		return ErrNotExist
	}

	return nil
}
