package router

import (
	"sync"
)

var once sync.Once

// Init init router
func Init() {
	once.Do(func() {
		initDNSHandler()
		initHTTPHandler()
	})
}
