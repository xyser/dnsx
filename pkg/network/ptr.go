package network

import (
	"net"
)

// DomainSuffixLen ptr domain suffix length
const DomainSuffixLen = 14

// PTRToIP 提取 PTR 的 IP
// 4.3.2.1.in-addr.arpa. to 1.2.3.4
func PTRToIP(domain []byte) string {
	// 移除后面 14个 字符
	ip := domain[0 : len(domain)-DomainSuffixLen]
	reverseWords(ip)
	return net.ParseIP(string(ip)).String()
}

// reverseWords reverse words
func reverseWords(s []byte) {
	l := len(s)
	reverse(s, 0, l-1)
	reverseWord(s, l)
}

// reverseWord reverse word
func reverseWord(s []byte, n int) {
	i, j := 0, 0

	for i < n {
		for i < j || (i < n && s[i] == '.') {
			i++
		}
		for j < i || (j < n && s[j] != '.') {
			j++
		}
		reverse(s, i, j-1)
	}
}

// reverse reverse bytes
func reverse(s []byte, i, j int) {
	if len(s) == 0 {
		return
	}
	var t byte
	for i < j {
		t = s[i]
		s[i] = s[j]
		s[j] = t
		i++
		j--
	}
}
