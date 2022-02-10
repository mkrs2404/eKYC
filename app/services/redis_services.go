package services

import (
	"context"
	"time"

	"github.com/mkrs2404/eKYC/app/redis_client"
)

//SetToRedis sets the value against the key in Redis
func SetToRedis(key string, value interface{}, expiration time.Duration) error {

	var redisClient = redis_client.GetRedisClient()
	err := redisClient.Set(context.Background(), key, value, expiration).Err()
	return err
}

//GetFromRedis returns the string value for the given key
func GetFromRedis(key string) (string, error) {

	var redisClient = redis_client.GetRedisClient()
	value, err := redisClient.Get(context.Background(), key).Result()
	return value, err
}
