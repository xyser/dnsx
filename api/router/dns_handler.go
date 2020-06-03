package router

import (
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type DNSCall func(r *dns.Msg) error

type Engine struct {
	handles sync.Map
}

func New() (h *Engine) {
	return &Engine{handles: sync.Map{}}
}

func (h *Engine) Register(qtype uint16, handle DNSCall) {
	h.handles.Store(qtype, handle)
}

func (h *Engine) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	//w.RemoteAddr()
	msg := dns.Msg{}
	msg.Authoritative = true       // 是否权威服务
	msg.RecursionAvailable = false // 是否递归查询响应

	msg.SetReply(r)
	if len(r.Question) == 0 {
		msg.Rcode = dns.RcodeFormatError
		_ = w.WriteMsg(r)
		return
	}

	if call, ok := h.handles.Load(r.Question[0].Qtype); ok {
		if callFunc, ok := call.(DNSCall); ok {
			if err := callFunc(&msg); err != nil {
				msg.Rcode = dns.RcodeServerFailure
				_ = w.WriteMsg(r)
				return
			}
		}
	}

	// 签名
	if len(msg.Answer) > 0 {
		_ = w.WriteMsg(r)
		return
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
