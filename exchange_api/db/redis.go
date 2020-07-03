package db

import (
	"exchange_api/cfg"

	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client {
	redisConf := cfg.GetStringMap("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf["host"].(string) + ":" + redisConf["port"].(string),
		Password: redisConf["auth"].(string), // no password set
		DB:       redisConf["db"].(int),      // use default DB
	})

	return client
}
