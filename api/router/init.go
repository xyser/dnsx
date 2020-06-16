package router

import (
	"sync"
)

var once sync.Once

func Init() {
	once.Do(func() {
		initDNSHandler()
		initHTTPHandler()
	})
}
