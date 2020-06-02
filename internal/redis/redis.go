package redis

import (
	"fmt"
	"sync"

	"dnsx/pkg/config"

	"github.com/go-redis/redis"
)

var once sync.Once
var client *redis.Client

const Nil = redis.Nil

func Init() {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     config.GetString("redis.addr"),
			Password: config.GetString("redis.password"), // no password set
			DB:       config.GetInt("redis.db"),          // use default DB
			PoolSize: config.GetInt("redis.pool"),
		})

		pong, err := client.Ping().Result()
		if err == nil {
			fmt.Printf("\033[1;30;42m[info]\033[0m redis connect success %s\n", pong)
		} else {
			panic(fmt.Sprintf("\033[1;30;41m[error]\033[0m redis connect error %s\n", err.Error()))
		}
	})
}

func Redis() *redis.Client {
	return client
}
