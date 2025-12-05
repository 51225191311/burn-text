package storage

import (
	"burn-text/internal/config"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background() //需要一个context上下文

func InitRedis() error {
	//从全局配置读取Redis信息
	cfg := config.GlobalConfig.Redis

	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr, //读取配置
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis连接失败（addr: %s): %v", cfg.Addr, err)
	}
	return nil
}

func Save(id string, data string, duration time.Duration) error {
	return rdb.Set(ctx, id, data, duration).Err()
}

func GetAndDelete(id string) (string, error) {
	val, err := rdb.Get(ctx, id).Result() //获取数据

	if err == redis.Nil {
		return "", errors.New("信息不存在或已过期")
	} else if err != nil {
		return "", err
	}

	rdb.Del(ctx, id)

	return val, nil
}

func AllowRequest(ip string, limit int64, window time.Duration) (bool, error) {
	//构建专属key
	key := "rate_limit:" + ip

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rdb.Expire(ctx, key, window)
	}

	if count > limit {
		return false, nil
	}

	return true, nil
}
