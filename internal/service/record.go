package service

import (
	"github.com/miekg/dns"

	v1 "github.com/dingdayu/dnsx/api/controller/v1"
	"github.com/dingdayu/dnsx/model/dao"
)

// GetRecordList get record list api by service
func GetRecordList(name, qtype, value string) (rrs []dao.Record, err error) {
	query := dao.Record{}
	if len(name) > 0 {
		query.Name = dns.Fqdn(name)
	}
	if len(qtype) > 0 {
		query.Type = qtype
	}
	if len(value) > 0 {
		query.Value = value
	}
	rrs, err = dao.GetRecord(query)
	if err != nil {
		return rrs, v1.ErrNotExist
	}
	return
}

// DefaultRecordTTL record default ttl
const DefaultRecordTTL = 3600

// CreateRecord create record
func CreateRecord(name, qtype, value string, ttl uint32, priority int) (rr dao.Record, err error) {
	if ttl == 0 {
		ttl = DefaultRecordTTL
	}

	rr = dao.Record{
		Name:     dns.Fqdn(name),
		Type:     qtype,
		Value:    value,
		TTL:      ttl,
		Priority: priority,
	}

	if err := dao.CreateRecord(&rr); err != nil {
		return rr, v1.ErrInternalServer
	}
	return
}
