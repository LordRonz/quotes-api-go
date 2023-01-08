package redisclient

import (
	"backend-2/api/cmd/utils"
	"context"

	"github.com/go-redis/redis/v8"
)

func GetClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.GetEnv("REDIS_URL", "localhost:6379"),
		Password: utils.GetEnv("REDIS_PASS", ""), // no password set
		DB:       0,                              // use default DB
	})
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

var Rdb = GetClient()
