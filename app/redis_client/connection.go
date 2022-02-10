package redis_client

import (
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitializeRedis(addr, password string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:        addr,
		IdleTimeout: 2 * time.Minute,
		PoolSize:    5,
		Password:    password,
	})
}
func GetRedisClient() *redis.Client {
	return redisClient
}
