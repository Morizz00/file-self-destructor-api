package storage

import (
	"context"
	"encoding/json"
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

func StoreFile(key string, file StoredFile) error {
	u, err := json.Marshal(file)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, u, time.Minute*5).Err()
}

func GetAndDelete(key string) (StoredFile, error) {
	val, err := rdb.GetDel(ctx, key).Bytes()
	if err != nil {
		return StoredFile{}, err
	}
	var res StoredFile
	err = json.Unmarshal(val, &res)
	return res, err
}
