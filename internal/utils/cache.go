package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InvalidateCache(ctx context.Context, client *redis.Client, pattern string) error {
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
