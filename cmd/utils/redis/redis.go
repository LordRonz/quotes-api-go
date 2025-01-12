package redisclient

import (
	"backend-2/api/cmd/utils"
	"context"
	"crypto/tls"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/go-redis/redis/v8"
)

func SetClient() *redis.Client {
	redisUrl := utils.GetEnv("REDIS_URL", "localhost:6379")
	if len(strings.Split(redisUrl, ":")) == 1 {
		redisUrl += ":6379"
	}
	redisPass := utils.GetEnv("REDIS_PASS", "")
	rdb := redis.NewClient(&redis.Options{
		Addr:      redisUrl,
		Password:  redisPass, // no password set
		DB:        0,         // use default DB
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	log.Printf("Connected to Redis at %v", redisUrl)

	Rdb = rdb
	return rdb
}

func DelByPattern(pattern string) {
	var cursor uint64
	keys, _, err := Rdb.Scan(Ctx, cursor, pattern, 0).Result()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		Rdb.Del(Ctx, key)
	}
}

func Clear() {
	DelByPattern("*")
}

var Ctx = context.Background()

var Rdb *redis.Client
