package database

import (
	"github.com/cheekybits/genny/generic"
	"github.com/go-redis/redis"
	"os"
)

var client = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: os.Getenv("REDIS_PASS"), // no password set
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

func HSetKey(key,field string,value generic.Type){
	_ = client.HSetNX(key,field,value)
}

func HSetIncrKey(key,field string){
	_ = client.HIncrBy(key,field,1)
}

func HGetKey(key,field string) string {
	return client.HGet(key,field).Val()
}

func IncrValue(key string) {
	_ = client.Incr(key)
}

func CountKeys (key string) int {
	keys := client.Keys(key).Val()
	return len(keys)
}

func HGetAllKey(key string) map[string]string {
	return client.HGetAll(key).Val()
}