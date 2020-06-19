package tests

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/panjf2000/ants/v2"
)

var rdb *redis.Client

func TestMain(m *testing.M) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "00000000", // no password set
		DB:       3,          // use default DB
		PoolSize: 100,
	})

	if pong, err := rdb.Ping().Result(); err != nil {
		fmt.Println("redis connect error:", err)
		return
	} else {
		fmt.Println(pong, err)
	}

	m.Run()
	fmt.Println("end")
}

func TestAnts(t *testing.T) {
	// Use the common pool.
	var wg sync.WaitGroup

	an, _ := ants.NewPool(500)

	// Submit tasks one by one.
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		an.Submit(func() {
			start := time.Now()
			if _, err := rdb.SetNX(strconv.Itoa(i), "value", 10*time.Second).Result(); err != nil {
				t.Error(err)
			}
			if latency := time.Since(start); latency.Microseconds() > 50 {
				t.Log(latency)
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
