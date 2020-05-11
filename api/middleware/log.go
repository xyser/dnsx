package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"dnsx/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WriterLog 处理跨域请求,支持options访问
func WriterLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)

		body := bodyBuf.Bytes()
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		latency := time.Since(start)
		if c.Request.Method == http.MethodGet {
			log.New().Named("http_request").Info(c.Request.RequestURI,
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Int("status", c.Writer.Status()),
				zap.String("ip", c.ClientIP()),
				zap.Int64("latency", latency.Nanoseconds()/1e6),
				zap.Any("response", blw.body.String()))
		} else {
			log.New().Named("http_request").Info(c.Request.RequestURI,
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Int("status", c.Writer.Status()),
				zap.String("ip", c.ClientIP()),
				zap.Int64("latency", latency.Nanoseconds()/1e6),
				zap.Any("response", blw.body.String()),
				zap.ByteString("body", body))
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
