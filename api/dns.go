package api

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/dingdayu/dnsx/model/dao"
	"github.com/dingdayu/dnsx/pkg/network"

	"github.com/miekg/dns"
)

type DNSHandler struct{}

func (h *DNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	//w.RemoteAddr()
	msg := dns.Msg{}
	msg.Authoritative = true       // 是否权威服务
	msg.RecursionAvailable = false // 是否递归查询响应

	msg.SetReply(r)
	if len(r.Question) == 0 {
		return
	}

	switch r.Question[0].Qtype {
	// 域名解析 IPV4,IPV6
	case dns.TypeA:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				A:   net.ParseIP(rr.Value),
			})
		}
	case dns.TypeAAAA:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				AAAA: net.ParseIP(rr.Value),
			})
		}
	// IP 解析域名
	case dns.TypePTR:
		msg.Authoritative = true
		ip := network.PTRToIP([]byte(msg.Question[0].Name))

		rrs, err := dao.GetRecord(map[string]interface{}{"type": "a", "value": ip})
		if err != nil {
			break
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.PTR{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Ptr: rr.Name,
			})
		}
	case dns.TypeCNAME:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.CNAME{
				Hdr:    dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Target: rr.Value,
			})
		}
	case dns.TypeNS:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.NS{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Ns:  rr.Value,
			})
		}
	case dns.TypeMX:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.MX{
				Hdr:        dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Preference: uint16(rr.Priority),
				Mx:         rr.Value,
			})
		}
	case dns.TypeSRV:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			// 优先级 权重 端口 目标地址 1 1 80 www.baidu.com
			vs := strings.Split(rr.Value, " ")
			var priority, weight, port uint64
			if len(vs) == 4 {
				priority, _ = strconv.ParseUint(vs[0], 16, 16)
				weight, _ = strconv.ParseUint(vs[1], 16, 16)
				port, _ = strconv.ParseUint(vs[2], 16, 16)
			} else if len(vs) == 3 {
				priority = uint64(rr.Priority)
				weight, _ = strconv.ParseUint(vs[0], 16, 16)
				port, _ = strconv.ParseUint(vs[1], 16, 16)
			}
			msg.Answer = append(msg.Answer, &dns.SRV{
				Hdr:      dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Priority: uint16(priority), // 优先级
				Weight:   uint16(weight),   // 权重
				Port:     uint16(port),     // 端口
				Target:   vs[len(vs)-1],    // 对应目标地址,可以是域名或IP
			})
		}
	case dns.TypeURI:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			// 优先级 权重 端口 目标地址 1 1 80 www.baidu.com
			vs := strings.Split(rr.Value, " ")
			var priority, weight uint64
			if len(vs) == 3 {
				priority, _ = strconv.ParseUint(vs[0], 16, 16)
				weight, _ = strconv.ParseUint(vs[1], 16, 16)
			} else if len(vs) == 2 {
				priority = uint64(rr.Priority)
				weight, _ = strconv.ParseUint(vs[0], 16, 16)
			}

			// URI 协议解释: https://www.dynu.com/Resources/DNS-Records/URI-Record
			msg.Answer = append(msg.Answer, &dns.URI{
				Hdr:      dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Priority: uint16(priority), // 优先级
				Weight:   uint16(weight),   // 权重
				Target:   vs[len(vs)-1],    // 对应目标地址,可以是域名或IP
			})
		}
	case dns.TypeTXT:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.TXT{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Txt: []string{rr.Value},
			})
		}
	case dns.TypeCAA:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		// CAA 协议解释: https://support.dnsimple.com/articles/caa-record/
		// RFC 文档: https://tools.ietf.org/html/rfc6844#section-3
		for _, rr := range rrs {
			// 标志位[0,1] 标签位[issue,issuewild,iodef] 机构域
			// domain.com. CAA 0 iodef mailto:admin@domain.com
			vs := strings.Split(rr.Value, " ")
			// 不够3位 不返回
			if len(vs) < 3 {
				continue
			}

			flag := 0
			if vs[0] == "1" {
				flag = 1
			}
			msg.Answer = append(msg.Answer, &dns.CAA{
				Hdr:  dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Flag: uint8(flag), // 标志位，严格校验 Tag 标签位
				// Tag 标签位
				// issue: 显式地授权单个证书颁发机构为主机名颁发证书（任何类型）。
				// issuewild: 显式地授权单个证书颁发机构为主机名颁发通配符证书（并且仅通配符）。
				// iodef: 指定证书颁发机构可以向其报告策略违规的URL。使用了事件对象描述交换格式（IODEF）格式
				Tag:   vs[1], // 标签位，
				Value: vs[2],
			})
		}
	case dns.TypeTLSA: // TLSA记录格式: 保存证书关联数据
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		// TLSA 协议解释: https://www.dynu.com/Resources/DNS-Records/TLSA-Record
		for _, rr := range rrs {
			vs := strings.Split(rr.Value, " ")
			// 不够3位 不返回
			if len(vs) < 4 {
				continue
			}
			usage, _ := strconv.ParseUint(vs[0], 8, 8)
			selector, _ := strconv.ParseUint(vs[1], 8, 8)
			match, _ := strconv.ParseUint(vs[2], 8, 8)

			msg.Answer = append(msg.Answer, &dns.TLSA{
				Hdr:          dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Usage:        uint8(usage),    // 证书使用情况
				Selector:     uint8(selector), // 选择器,0-完整证书,1-使用主题公钥
				MatchingType: uint8(match),    // 匹配类型, 0-无哈希, 1-所选内容的SHA-256哈希,2-所选内容的SHA-512哈希
				Certificate:  "",              // 证书关联数据
			})
		}
	case dns.TypeHINFO:
		domain := msg.Question[0].Name
		rrs, err := dao.GetRecordByNameAndType(domain, r.Question[0].Qtype)
		if err != nil {
			break
		}
		for _, rr := range rrs {
			vs := strings.Split(rr.Value, " ")
			// 不够3位 不返回
			if len(vs) < 2 {
				continue
			}
			msg.Answer = append(msg.Answer, &dns.HINFO{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				Cpu: vs[0],
				Os:  vs[1],
			})
		}
	}

	// DNSSEC https://support.cloudflare.com/hc/zh-cn/articles/360006660072
	// DNSSEC https://backreference.org/2010/11/17/dnssec-verification-with-dig/
	// DNSSEC https://dnsviz.net/d/git.xyser.net/dnssec/
	// ADDITIONAL https://www.jianshu.com/p/71f61652ec23
	// 我所理解的 DNSSEC https://imlonghao.com/41.html
	// 巧妙运用DNS及其安全扩展DNSSec https://zhuanlan.zhihu.com/p/52877648
	// DNSSEC的概念及作用 https://www.cloudxns.net/Support/detail/id/1309.html
	// https://tools.ietf.org/html/rfc6781

	// 查询上游服务器
	if len(msg.Answer) == 0 && msg.RecursionDesired {
		qr, _, _ := QuestionStream(r.Question[0].Name, r.Question[0].Qtype)
		if r == nil {
			dns.HandleFailed(w, r)
			return
		}

		if qr.Rcode == dns.RcodeSuccess {
			msg.Answer = qr.Answer
			if len(qr.Ns) > 0 {
				msg.Ns = qr.Ns
			}
			if len(qr.Extra) > 0 {
				msg.Extra = qr.Extra
			}
		}
		// 本次查询非权威应答
		msg.Authoritative = false
		msg.RecursionAvailable = true // 是否递归查询响应
	}

	_ = w.WriteMsg(&msg)
}

// QuestionStream 询问上游服务器
func QuestionStream(name string, qtype uint16) (r *dns.Msg, rtt time.Duration, err error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	// 启用 EDNS
	m.SetEdns0(4096, true)
	m.AuthenticatedData = true
	// 设置递归查询
	m.SetQuestion(dns.Fqdn(name), qtype)
	return c.Exchange(m, net.JoinHostPort("8.8.8.8", "53"))
}
