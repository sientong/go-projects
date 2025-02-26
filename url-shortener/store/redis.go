package store

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func SaveURL(ctx context.Context, shortURL, longURL string) error {
	return rdb.Set(ctx, shortURL, longURL, 0).Err()
}

func GetURL(ctx context.Context, shortURL string) (string, error) {
	return rdb.Get(ctx, shortURL).Result()
}
