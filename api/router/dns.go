package router

import (
	"github.com/miekg/dns"

	dnsc "github.com/dingdayu/dnsx/api/controller/dns"
)

// DNSHandler DNS 路由注册
// https://zh.wikipedia.org/zh/DNS%E8%AE%B0%E5%BD%95%E7%B1%BB%E5%9E%8B%E5%88%97%E8%A1%A8
// https://www.dynu.com/Resources/DNS-Records
func DNSHandler() *Engine {
	handle := New()

	handle.Register(dns.TypeA, dnsc.TypeA)
	handle.Register(dns.TypeAAAA, dnsc.TypeAAAA)

	// IP 解析域名
	handle.Register(dns.TypePTR, dnsc.TypePTR)
	handle.Register(dns.TypeCNAME, dnsc.TypeCNAME)

	// TypeNS 名称服务器记录
	handle.Register(dns.TypeNS, dnsc.TypeNS)

	// TypeMX 电邮交互记录
	handle.Register(dns.TypeMX, dnsc.TypeMX)

	// TypeSRV Service locator
	// _sip._udp 3600 IN SRV 10 5 5060 siphost.com.
	handle.Register(dns.TypeSRV, dnsc.TypeSRV)

	// TypeURI SRV to URI
	// _ftp._tcp.example.com. 3600 IN URI 10 1 "ftp://ftp.example.com/public"
	handle.Register(dns.TypeURI, dnsc.TypeURI)

	handle.Register(dns.TypeTXT, dnsc.TypeTXT)

	// TypeCAA
	handle.Register(dns.TypeCAA, dnsc.TypeCAA)

	// TypeTLSA TLSA记录格式: 保存证书关联数据
	handle.Register(dns.TypeTLSA, dnsc.TypeTLSA)

	// TypeHINFO hardware type and operating system (OS) information
	// testhinfo 90 IN HINFO "INTEL-386" "Windows"
	handle.Register(dns.TypeHINFO, dnsc.TypeHINFO)

	return handle
}
