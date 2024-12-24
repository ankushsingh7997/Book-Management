package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/ankush/bookstore/env"
	"github.com/ankush/bookstore/logger"
	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	Client       *redis.Client
	Connected    bool
	LastErrorMsg string
}

func prepareRedisConfig() *redis.Options {
	redisPassword := env.Get("REDIS_PASSWORD", "")
	redisHost := env.Get("REDIS_HOST", "localhost")
	redisPort := env.Get("REDIS_PORT", "6379")
	redisName := env.Get("REDIS_NAME", "default")
	var redisURL string
	if redisPassword == "" || redisPassword == "NONE" {
		redisURL = fmt.Sprintf("redis://%s:%s", redisHost, redisPort)
	} else {
		redisURL = fmt.Sprintf("redis://:%s@%s:%s", redisPassword, redisHost, redisPort)
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Printf("Error parsing Redis URL: %v", err)
	}

	opts.DialTimeout = 10 * time.Second // Connection timeout
	opts.MaxRetries = 3                 // Retry on failure
	opts.Username = redisName           // Set Redis username (if applicable)

	return opts
}

func initializeRedisConnection() *RedisConnection {
	logg := logger.NewLogger("UserService", "production")
	rc := &RedisConnection{
		Connected: false,
	}

	opts := prepareRedisConfig()
	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		rc.Connected = false
		rc.LastErrorMsg = err.Error()
		logg.Error(fmt.Sprintf("Redis Connection Error: %v", err))
		return rc
	}

	rc.Client = client
	rc.Connected = true

	logg.Info(fmt.Sprintf("Redis client connected to %s", opts.Username))

	return rc
}

// func createRetryStrategy(retries int, cause error) time.Duration {
// 	fmt.Println("Reconnection attempt:%d", retries)
// 	if cause == nil {
// 		fmt.Println("Retry cause: %v", cause)
// 	}
// 	return time.Duration(min(retries*1000, 120000)) * time.Millisecond
// }
// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

var RedisConn *RedisConnection

func init() {
	RedisConn = initializeRedisConnection()
}
