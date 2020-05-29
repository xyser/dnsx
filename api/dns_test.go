package api

import (
	"fmt"
	"net"
	"testing"

	"github.com/miekg/dns"
)

func TestUpStream(t *testing.T) {
	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("qq.com"), dns.TypeA)
	//m.SetQuestion(dns.Fqdn("taobao.com"), dns.TypeHINFO)

	m.Question = append(m.Question, dns.Question{
		Name:   dns.Fqdn("baidu.com"),
		Qtype:  dns.TypeA,
		Qclass: 1,
	})
	m.Question = append(m.Question, dns.Question{
		Name:   dns.Fqdn("taobao.com"),
		Qtype:  dns.TypeA,
		Qclass: 1,
	})
	r, _, err := c.Exchange(m, net.JoinHostPort("1.1.1.1", "53"))
	fmt.Println(r, err)

	fmt.Printf("%+v", m.Question)
}
