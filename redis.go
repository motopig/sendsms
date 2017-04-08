package smservice

import (
	"time"

	"strconv"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func storeCode(mobile string, code string) bool {

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	err = redisClient.Set(mobile, code, time.Minute*5).Err()
	if err != nil {
		return false
	} else {
		return true
	}
}

func getCode(mobile string) string {

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	rcode := redisClient.Get(mobile).Val()
	if rcode != "" {
		return rcode
	} else {
		return `none`
	}
}

func initRedis() *redis.Client {
	db, _ := strconv.Atoi(config.RedisConf["db"])
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisConf["host"] + ":" + config.RedisConf["port"],
		Password: config.RedisConf["password"], // no password set
		DB:       db,                           // use default DB
	})
}
