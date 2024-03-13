package store

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/youth-service/auth/pkg/setting"
	"strconv"
	"time"
)

var rdb *redis.Client

func RunRedis() {

	rdb = redis.NewClient(&redis.Options{
		Addr:            setting.RedisSetting.Host,
		Password:        setting.RedisSetting.Password, // no password set
		DB:              setting.RedisSetting.DB,       // use default DB
		PoolSize:        setting.RedisSetting.PoolSize,
		MinIdleConns:    setting.RedisSetting.MinIdleConns,
		ConnMaxIdleTime: setting.RedisSetting.ConnMaxIdleTime,
	})

}
func Get(key string) ([]byte, error) {
	reply, err := rdb.Get(context.TODO(), key).Bytes()
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func SaveMetadata(ctx context.Context, uuid string, uid string, etime int64) error {
	expiresTime := time.Unix(etime, 0)
	now := time.Now()

	err := rdb.Set(ctx, uuid, uid, expiresTime.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetPassportData(ctx context.Context, uuid string) (uint64, error) {
	result, err := rdb.Get(ctx, uuid).Result()
	if err != nil {
		return 0, err
	}

	uid, _ := strconv.ParseUint(result, 10, 64)
	return uid, nil
}

func DeleteMetadata(ctx context.Context, uuid string) (int64, error) {
	deleted, err := rdb.Del(ctx, uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
