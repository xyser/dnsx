package entity

import (
	"time"

	"github.com/miekg/dns"
)

type AnswerCache struct {
	Question dns.Question `json:"question"`
	Answer   []string     `json:"answer"`
	Extra    []string     `json:"extra"`
	Ns       []string     `json:"ns"`

	RecursionAvailable bool `json:"recursion_available"`
	Authoritative      bool `json:"authoritative"`

	Expire time.Time `json:"expire"`
}
