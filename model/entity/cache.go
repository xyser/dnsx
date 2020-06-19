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
func (a AnswerCache) ToMsg(msg *dns.Msg) {
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
