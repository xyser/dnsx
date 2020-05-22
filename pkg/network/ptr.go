package network

import (
	"net"
)

const DomainSuffixLen = 14

// PTRToIP 提取 PTR 的 IP
// 4.3.2.1.in-addr.arpa. to 1.2.3.4
func PTRToIP(domain []byte) string {
	// 移除后面 14个 字符
	ip := domain[0 : len(domain)-DomainSuffixLen]
	reverse(ip)
	return net.ParseIP(string(ip)).String()
}

// reverse 翻转字符串
func reverse(s []byte) {
	if len(s) == 0 {
		return
	}
	i := 0
	j := len(s) - 1
	var t byte
	for i < j {
		t = s[i]
		s[i] = s[j]
		s[j] = t
		i++
		j--
	}
}
