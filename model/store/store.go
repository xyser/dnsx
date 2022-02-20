package store

import (
	"github.com/miekg/dns"
)

// Store interface
type Store interface {
	// Question query
	Question(question dns.Question) []dns.RR

	// CreateDomain domain api
	CreateDomain(domain string)
	DeleteDomain(domain string)
	DomainList()

	// CreateRecord record api
	CreateRecord(domain string, question dns.Question, rr string)
	DeleteRecord()
	RecordList()
}
