package pprof

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/dingdayu/dnsx/pkg/config"
	"github.com/dingdayu/dnsx/pkg/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var mu sync.Mutex
var ctx context.Context
var server *http.Server
var isOpen bool

const LogNamed = "pprof"

// Init init server
func Init() {
	config.RegisterChangeEvent(func(e fsnotify.Event) {
		ReloadConfig()
	})
}

func ReloadConfig() {
	// 配置未开启状态(false)
	if !viper.GetBool("pprof.open") {
		// 服务已开启,执行关闭
		if getStatus() {
			setStatus(false)
			go server.Shutdown(ctx)
		}
		return
	}

	// 配置开启状态(true)且服务未开启
	if !getStatus() {
		ctx = context.TODO()
		server = &http.Server{
			Addr: viper.GetString("pprof.addr"),
		}
		setStatus(true)

		go server.ListenAndServe()
		log.New().Named(LogNamed).Info("start pprof", zap.String("url", viper.GetString("pprof.addr")))
		return
	}
}

// setStatus 设置 pprof 是否开启
func setStatus(status bool) {
	mu.Lock()
	defer mu.Unlock()
	isOpen = status
}

// getStatus 获取 pprof 是否开启
func getStatus() bool {
	mu.Lock()
	defer mu.Unlock()
	return isOpen
}
