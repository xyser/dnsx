package service

import (
	v1 "dnsx/api/controller/v1"
	"dnsx/model/dao"
	"github.com/miekg/dns"
)

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

// DefaultRecordTTL 记录默认TTL
const DefaultRecordTTL = 3600

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
