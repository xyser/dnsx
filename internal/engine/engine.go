package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xyser/dnsx/model/entity"
	"github.com/xyser/dnsx/pkg/config"

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

// cacheKey sum cache key
func (h *Engine) cacheKey(msg *dns.Msg) string {
	return fmt.Sprintf("%s:%d:%d", msg.Question[0].Name, msg.Question[0].Qtype, msg.Question[0].Qclass)
}

// StoreCache store record cache
func (h *Engine) StoreCache(msg *dns.Msg) (err error) {
	if len(msg.Answer) == 0 {
		return nil
	}
	var cah entity.AnswerCache
	cah.ToCache(msg)

	if mar, err := json.Marshal(cah); err == nil {
		return h.cache.Set(h.cacheKey(msg), mar)
	}
	return nil
}

// LoadCache load record cache
func (h *Engine) LoadCache(msg *dns.Msg) (hasCache bool, err error) {
	if b, err := h.cache.Get(h.cacheKey(msg)); err == nil {
		var cache entity.AnswerCache
		if err = json.Unmarshal(b, &cache); err == nil {
			// expire
			if cache.Expire.Before(time.Now()) {
				_ = h.cache.Delete(h.cacheKey(msg))
				return false, errors.New("nil")
			}
			if len(cache.Answer) > 0 {
				cache.ToMsg(msg)
				return true, nil
			}
		}
	}
	return false, errors.New("nil")
}

// Handle engine handle
func (h *Engine) Handle(msg *dns.Msg) (err error) {
	// load cache
	if has, _ := h.LoadCache(msg); has {
		return
	}

	// handle
	if call, ok := h.handles.Load(msg.Question[0].Qtype); ok {
		if callFunc, ok := call.(DNSCall); ok {
			if err := callFunc(msg); err != nil {
				msg.Rcode = dns.RcodeServerFailure
				return err
			}
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
		qr, _, err := QuestionStream(msg.Question[0].Name, msg.Question[0].Qtype)
		if err != nil {
			return err
		}

		if qr.Rcode == dns.RcodeSuccess {
			msg.Answer = qr.Answer
			msg.Ns = qr.Ns
			msg.Extra = qr.Extra
		}
		// 本次查询非权威应答
		msg.Authoritative = false
		msg.RecursionAvailable = true // 是否递归查询响应
	}

	// store cache
	_ = h.StoreCache(msg)
	return nil
}

// ErrNotConfigUpstream not config upstream
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
