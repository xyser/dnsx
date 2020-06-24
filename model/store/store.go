package store

import (
	"github.com/miekg/dns"
)

// Store store interface
type Store interface {
	// query
	Question(question dns.Question) []dns.RR

	// domain api
	CreateDomain(domain string)
	DeleteDomain(domain string)
	DomainList()

	// record api
	CreateRecord(domain string, question dns.Question, rr string)
	DeleteRecord()
	RecordList()
}
