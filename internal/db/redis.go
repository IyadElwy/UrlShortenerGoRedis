package db

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrNotFoundInRedis = errors.New("Not Found")

var ctx = context.Background()

type Redis struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) (*Redis, error) {
	r := &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
	val, err := r.client.Ping(ctx).Result()
	if err != nil && val != "PONG" {
		return nil, err
	}
	return r, nil
}

func (r *Redis) Store(key, value string) error {
	err := r.client.Set(ctx, key, value, redis.KeepTTL).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Retrieve(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFoundInRedis
	}
	return val, nil
}
