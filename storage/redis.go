package storage

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	// Clean up REDIS_URL if it contains redis-cli command prefix
	redisURL = cleanRedisURL(redisURL)

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic("Failed to parse REDIS_URL: " + err.Error() + " (got: " + redisURL + ")")
	}

	rdb = redis.NewClient(opt)
	
	// Test connection with timeout - don't block startup
	// If Redis is unavailable, operations will fail gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	if err := rdb.Ping(ctx).Err(); err != nil {
		// Log warning but don't panic - allows app to start
		// Redis operations will return errors that can be handled
		// This is important for deployment environments where Redis might start after the app
	}
}

// cleanRedisURL removes common redis-cli command prefixes
func cleanRedisURL(url string) string {
	// Remove "redis-cli --tls -u " prefix if present
	if strings.HasPrefix(url, "redis-cli") {
		// Extract URL after "-u " or "--tls -u "
		parts := strings.Split(url, "-u ")
		if len(parts) > 1 {
			url = strings.TrimSpace(parts[len(parts)-1])
		}
	}
	return strings.TrimSpace(url)
}

func StoreFile(key string, file StoredFile, expiry time.Duration) error {
	u, err := json.Marshal(file)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, u, expiry).Err()
}

//	func GetAndDelete(key string) (StoredFile, error) {
//		val, err := rdb.GetDel(ctx, key).Bytes()
//		if err != nil {
//			return StoredFile{}, err
//		}
//		var res StoredFile
//		err = json.Unmarshal(val, &res)
//		return res, err
//	}
func Get(key string) (StoredFile, error) {
	val, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return StoredFile{}, err
	}
	var res StoredFile
	err = json.Unmarshal(val, &res)
	return res, err
}
func Delete(key string) error {
	return rdb.Del(ctx, key).Err()
}
func UpdateFilePreservingTTL(key string, file StoredFile) error {
	u, err := json.Marshal(file)
	if err != nil {
		return err
	}
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		ttl = time.Minute * 5
	}
	return rdb.Set(ctx, key, u, ttl).Err()
}
