package database

import (
	"github.com/cheekybits/genny/generic"
	"github.com/go-redis/redis"
	"os"
)

var client = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDISCLOUD_URL"),
	Password: "", // no password set
	DB:       0,  // use default DB
})

func ConnectRedis() *redis.Client {
	return client
}

func SetKey(key string, value generic.Type) {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func GetKey(key string) string {
	val, err := client.Get(key).Result()
	if err != nil {
		return ""
	}
	return val
}
