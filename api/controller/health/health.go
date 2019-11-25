package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 首页，用于健康检查
func Hello(c *gin.Context) {
	c.String(http.StatusOK, "hello, word.")
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Prometheus 监控访问点
func Prometheus(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func init() {
	// 注册应用监控
	workerDB := NewClusterManager("db")
	prometheus.MustRegister(workerDB)
}
