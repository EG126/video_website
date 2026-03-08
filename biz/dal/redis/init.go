package redis

import (
	"context"
	"video_website/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Init() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})
	_, err := RDB.Ping(context.Background()).Result()
	return err
}
