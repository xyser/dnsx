package tests

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/panjf2000/ants/v2"
)

var rdb *redis.Client

const (
	// redis add and password
	redisAddr = "0.0.0.0:6379"
	redisPass = ""
)

// InitRedis init test redis
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass, // no password set
		DB:       3,         // use default DB
		PoolSize: 100,
	})

	if pong, err := rdb.Ping().Result(); err == nil {
		fmt.Println(pong, err)
	} else {
		fmt.Println("redis connect error:", err)
	}
}

func MockRedis() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(&redis.Options{Addr: mr.Addr()})

	if pong, err := rdb.Ping().Result(); err != nil {
		fmt.Println("redis connect error:", err)
	} else {
		fmt.Println(pong, err)
	}
}

func BenchmarkRedisSet(b *testing.B) {
	//InitRedis()
	MockRedis()

	// Use the common pool.
	var wg sync.WaitGroup

	an, _ := ants.NewPool(500)

	// Submit tasks one by one.
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		an.Submit(func() {
			start := time.Now()
			if _, err := rdb.SetNX(strconv.Itoa(i), "value", 10*time.Second).Result(); err != nil {
				b.Error(err)
			}
			if latency := time.Since(start); latency.Microseconds() > 50 {
				b.Log(latency)
			}
			wg.Done()
		})
	}
	wg.Wait()
}

func BenchmarkSetNX(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := rdb.SetNX("i", "value", 10*time.Second).Result(); err != nil {
			b.Error(err)
		}
	}
}
