package db

import (
	"exchange_api/config"
	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_HOSTS + ":" + config.REDIS_PORT,
		Password: config.REDIS_PASSWORD, // no password set
		DB:       config.REDIS_DB,  // use default DB
	})

	return client
}

