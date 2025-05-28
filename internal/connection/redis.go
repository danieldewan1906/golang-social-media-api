package connection

import (
	"context"
	"golang-restful-api/internal/config"
	"golang-restful-api/internal/util"

	"github.com/redis/go-redis/v9"
)

func GetRedisConnection(config config.Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		util.PanicIfError(err)
	}

	return rdb
}
