package router

import (
	"dnsx/api/controller/health"
	v1 "dnsx/api/controller/v1"
	"dnsx/api/controller/v1/record"
	"dnsx/api/middleware"
	"dnsx/pkg/config"

	"github.com/gin-gonic/gin"
)

var handle *gin.Engine

func Handler() *gin.Engine {
	handle = gin.New()
	// 正式环境不再在控制台请求输出日志
	if gin.Mode() != gin.ReleaseMode {
		handle.Use(gin.Logger())
	}
	handle.Use(gin.Recovery())

	handle.GET("/", health.Hello)
	handle.HEAD("/health", health.Hello)
	handle.GET("/health", health.Hello)
	handle.GET("/ping", health.Ping)
	handle.GET("/metrics", health.Prometheus)

	apiv1 := handle.Group("/api/v1")

	// 根据配置决定是否启用 api 请求日志
	if config.GetBool("log.request_log") {
		apiv1.Use(middleware.WriterLog())
	}
	apiv1.Use(middleware.Cors())

	// 服务路由
	apiv1.GET("/version", v1.Version)

	// domain
	apiv1.GET("/domains", v1.Version)
	apiv1.GET("/domains/:domain/records", v1.Version)

	// records
	apiv1.GET("/records", record.Lists)
	apiv1.POST("/records", record.Create)

	return handle
}
