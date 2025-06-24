package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func StoreFile(key string, data []byte) error {
	return rdb.Set(ctx, key, data, time.Minute*5).Err()
}

func GetAndDelete(key string) ([]byte, error) {
	val, err := rdb.GetDel(ctx, key).Bytes()
	return val, err
}
