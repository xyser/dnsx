package engine

import (
	"github.com/gin-gonic/gin"
)

type Option struct {
	tcpAddr string
	udpAddr string

	httpDNSAddr string
	httpAPIAddr string

	dnsHandle     *Engine
	httpAPIHandle *gin.Engine
	httpDNSHandle *gin.Engine
}

type OptionsFunc func(c *Option)

// WithDNSAddr option set tcp and udp Addr
func WithDNSAddr(addr string) OptionsFunc {
	return func(c *Option) {
		c.tcpAddr = addr
		c.udpAddr = addr
	}
}

// WithTCPAddr option set tcpAddr
func WithTCPAddr(tcpAddr string) OptionsFunc {
	return func(c *Option) {
		c.tcpAddr = tcpAddr
	}
}

// WithUDPAddr option set udpAddr
func WithUDPAddr(udpAddr string) OptionsFunc {
	return func(c *Option) {
		c.udpAddr = udpAddr
	}
}

// WithHttpDNSAddr option set httpDNSAddr
func WithHttpDNSAddr(httpDNSAddr string) OptionsFunc {
	return func(c *Option) {
		c.httpDNSAddr = httpDNSAddr
	}
}

// WithHttpAPIAddr option set httpAPIAddr
func WithHttpAPIAddr(httpAPIAddr string) OptionsFunc {
	return func(c *Option) {
		c.httpAPIAddr = httpAPIAddr
	}
}

// WithDNSHandle option set dnsHandle
func WithDNSHandle(dnsHandle *Engine) OptionsFunc {
	return func(c *Option) {
		c.dnsHandle = dnsHandle
	}
}

// WithHttpAPIHandle option set httpAPIHandle
func WithHttpAPIHandle(httpAPIHandle *gin.Engine) OptionsFunc {
	return func(c *Option) {
		c.httpAPIHandle = httpAPIHandle
	}
}

// WithHttpDNSHandle option set httpDNSHandle
func WithHttpDNSHandle(httpDNSHandle *gin.Engine) OptionsFunc {
	return func(c *Option) {
		c.httpDNSHandle = httpDNSHandle
	}
}
