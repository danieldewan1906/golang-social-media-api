package repository

import (
	"context"
	"encoding/json"
	"golang-restful-api/domains"
	"golang-restful-api/internal/util"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	rdb *redis.Client
}

var ctx = context.Background()

// Get implements domains.RedisRepository.
func (repository *redisRepository) Get(key string) (result string, err error) {
	data, _ := repository.rdb.Get(ctx, key).Result()
	return data, nil
}

// Set implements domains.RedisRepository.
func (repository *redisRepository) Set(key string, value interface{}) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		util.PanicIfError(err)
	}

	err = repository.rdb.Set(context.Background(), key, data, 5*time.Minute).Err()
	if err != nil {
		util.PanicIfError(err)
	}

	return err
}

func NewRedisRepository(redisConn *redis.Client) domains.RedisRepository {
	return &redisRepository{
		rdb: redisConn,
	}
}
