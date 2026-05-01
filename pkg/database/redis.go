package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisCache(cfg *config.Config) *redis.Client {
	host := cfg.RDBHost
	port := cfg.RDBPort
	addr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected")
	return rdb
}
