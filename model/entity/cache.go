package entity

import (
	"math"
	"time"

	"github.com/miekg/dns"
)

// AnswerCache answer cache
type AnswerCache struct {
	Question dns.Question `json:"question"`
	Answer   []string     `json:"answer"`
	Extra    []string     `json:"extra"`
	Ns       []string     `json:"ns"`

	RecursionAvailable bool `json:"recursion_available"`
	Authoritative      bool `json:"authoritative"`

	Expire time.Time `json:"expire"`
}

// ToMsg cache to Msg
func (a *AnswerCache) ToMsg(msg *dns.Msg) {
	var ttl float64
	if !a.Expire.IsZero() {
		ttl = a.Expire.Sub(time.Now()).Seconds()
	}

	for _, v := range a.Answer {
		if rr, err := dns.NewRR(v); err == nil {
			rr.Header().Ttl = uint32(math.Round(ttl))
			msg.Answer = append(msg.Answer, rr)
		}
	}
	for _, v := range a.Extra {
		if rr, err := dns.NewRR(v); err == nil && rr != nil {
			rr.Header().Ttl = uint32(math.Round(ttl))
			msg.Extra = append(msg.Extra, rr)
		}
	}
	for _, v := range a.Extra {
		if rr, err := dns.NewRR(v); err == nil && rr != nil {
			rr.Header().Ttl = uint32(math.Round(ttl))
			msg.Extra = append(msg.Extra, rr)
		}
	}
	msg.RecursionAvailable = a.RecursionAvailable
	msg.Authoritative = a.Authoritative
	return
}

// ToCache cache to Msg
func (a *AnswerCache) ToCache(msg *dns.Msg) {
	a.Question = msg.Question[0]
	a.Authoritative = msg.Authoritative
	a.RecursionAvailable = msg.RecursionAvailable
	a.Expire = time.Now().Add(time.Duration(msg.Answer[0].Header().Ttl) * time.Second)

	for _, v := range msg.Answer {
		a.Answer = append(a.Answer, v.String())
	}
	for _, v := range msg.Extra {
		a.Extra = append(a.Extra, v.String())
	}
	for _, v := range msg.Ns {
		a.Ns = append(a.Ns, v.String())
	}
	return
}
