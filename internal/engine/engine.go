package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dingdayu/dnsx/model/entity"
	"github.com/dingdayu/dnsx/pkg/config"

	"github.com/allegro/bigcache"
	"github.com/miekg/dns"
)

// DNSCall dns controller func
type DNSCall func(r *dns.Msg) error

// Engine dns engine
type Engine struct {
	handles sync.Map
	cache   *bigcache.BigCache
}

// NewEngine New new dns engine
func NewEngine() (h *Engine) {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	return &Engine{handles: sync.Map{}, cache: cache}
}

// Register register dns controller
func (h *Engine) Register(qtype uint16, handle DNSCall) {
	h.handles.Store(qtype, handle)
}

// ServeDNS export func to serve
func (h *Engine) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	//w.RemoteAddr()
	msg := dns.Msg{}
	msg.Authoritative = true       // 是否权威服务
	msg.RecursionAvailable = false // 是否递归查询响应

	msg.SetReply(r)
	if len(r.Question) == 0 {
		msg.Rcode = dns.RcodeFormatError
		_ = w.WriteMsg(&msg)
		return
	}

	if err := h.Handle(&msg); err != nil {
		msg.Rcode = dns.RcodeServerFailure
		_ = w.WriteMsg(&msg)
		return
	}
	_ = w.WriteMsg(&msg)
	return
}

// Handle engine handle
func (h *Engine) Handle(msg *dns.Msg) (err error) {
	// read cache
	ck := fmt.Sprintf("%s:%d:%d", msg.Question[0].Name, msg.Question[0].Qtype, msg.Question[0].Qclass)
	if b, err := h.cache.Get(ck); err == nil {
		temp := entity.AnswerCache{}
		if err = json.Unmarshal(b, &temp); err == nil {
			// expire
			if temp.Expire.Before(time.Now()) {
				_ = h.cache.Delete(ck)
				goto Handle
			}
			if len(temp.Answer) > 0 {
				ttl := temp.Expire.Sub(time.Now()).Seconds()
				for _, v := range temp.Answer {
					if rr, err := dns.NewRR(v); err == nil {
						rr.Header().Ttl = uint32(math.Round(ttl))
						msg.Answer = append(msg.Answer, rr)
					}
				}
				for _, v := range temp.Extra {
					if rr, err := dns.NewRR(v); err == nil && rr != nil {
						rr.Header().Ttl = uint32(math.Round(ttl))
						msg.Extra = append(msg.Extra, rr)
					}
				}
				for _, v := range temp.Extra {
					if rr, err := dns.NewRR(v); err == nil && rr != nil {
						rr.Header().Ttl = uint32(math.Round(ttl))
						msg.Extra = append(msg.Extra, rr)
					}
				}
				msg.RecursionAvailable = temp.RecursionAvailable
				msg.Authoritative = temp.Authoritative
				return err
			}
		}
	}

Handle:
	if call, ok := h.handles.Load(msg.Question[0].Qtype); ok {
		if callFunc, ok := call.(DNSCall); ok {
			if err := callFunc(msg); err != nil {
				msg.Rcode = dns.RcodeServerFailure
				return err
			}
		}
	}

	// 签名
	if len(msg.Answer) > 0 {
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
		qr, _, err := QuestionStream(msg.Question[0].Name, msg.Question[0].Qtype)
		if err != nil {
			return err
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

	// store cache
	if len(msg.Answer) > 0 {
		cah := entity.AnswerCache{
			Question:           msg.Question[0],
			Authoritative:      msg.Authoritative,
			RecursionAvailable: msg.RecursionAvailable,
			Expire:             time.Now().Add(time.Duration(msg.Answer[0].Header().Ttl) * time.Second),
		}
		for _, v := range msg.Answer {
			cah.Answer = append(cah.Answer, v.String())
		}
		for _, v := range msg.Extra {
			cah.Extra = append(cah.Extra, v.String())
		}
		for _, v := range msg.Ns {
			cah.Ns = append(cah.Ns, v.String())
		}
		mar, _ := json.Marshal(cah)
		_ = h.cache.Set(fmt.Sprintf("%s:%d:%d", msg.Question[0].Name, msg.Question[0].Qtype, msg.Question[0].Qclass), mar)
	}
	return nil
}

var ErrNotConfigUpstream = errors.New("not configuration upstream")

// QuestionStream query upstream
func QuestionStream(name string, qtype uint16) (r *dns.Msg, rtt time.Duration, err error) {
	// load upstream config
	upStream := config.GetString("app.upstream")
	if len(upStream) == 0 {
		return nil, 0, ErrNotConfigUpstream
	}

	// request upstream
	c := new(dns.Client)
	m := new(dns.Msg)
	// enable EDNS
	m.SetEdns0(4096, true)

	m.AuthenticatedData = true  // enable auth
	m.RecursionAvailable = true // enable recursive

	m.SetQuestion(dns.Fqdn(name), qtype)
	return c.Exchange(m, upStream)
}
