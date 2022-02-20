package router

import (
	dnsc "github.com/xyser/dnsx/api/controller/dns"
	"github.com/xyser/dnsx/internal/engine"

	"github.com/miekg/dns"
)

// dnsHandler dns handler
var dnsHandler *engine.Engine

// initDNSHandler init dns handler
func initDNSHandler() {
	dnsHandler = engine.NewEngine()

	dnsHandler.Register(dns.TypeA, dnsc.TypeA)
	dnsHandler.Register(dns.TypeAAAA, dnsc.TypeAAAA)

	// IP 解析域名
	dnsHandler.Register(dns.TypePTR, dnsc.TypePTR)
	dnsHandler.Register(dns.TypeCNAME, dnsc.TypeCNAME)

	// TypeNS 名称服务器记录
	dnsHandler.Register(dns.TypeNS, dnsc.TypeNS)

	// TypeMX 电邮交互记录
	dnsHandler.Register(dns.TypeMX, dnsc.TypeMX)

	// TypeSRV Service locator
	// Sample: _sip._udp 3600 IN SRV 10 5 5060 siphost.com.
	dnsHandler.Register(dns.TypeSRV, dnsc.TypeSRV)

	// TypeURI SRV to URI
	// Sample: _ftp._tcp.example.com. 3600 IN URI 10 1 "ftp://ftp.example.com/public"
	dnsHandler.Register(dns.TypeURI, dnsc.TypeURI)

	dnsHandler.Register(dns.TypeTXT, dnsc.TypeTXT)

	// TypeCAA
	dnsHandler.Register(dns.TypeCAA, dnsc.TypeCAA)

	// TypeTLSA TLSA记录格式: 保存证书关联数据
	dnsHandler.Register(dns.TypeTLSA, dnsc.TypeTLSA)

	// TypeHINFO hardware type and operating system (OS) information
	// Sample: testhinfo 90 IN HINFO "INTEL-386" "Windows"
	dnsHandler.Register(dns.TypeHINFO, dnsc.TypeHINFO)
}

// DNSHandler DNS 路由注册
// https://zh.wikipedia.org/zh/DNS%E8%AE%B0%E5%BD%95%E7%B1%BB%E5%9E%8B%E5%88%97%E8%A1%A8
// https://www.dynu.com/Resources/DNS-Records
func DNSHandler() *engine.Engine {
	return dnsHandler
}
