package dns

import (
	"github.com/miekg/dns"

	"github.com/dingdayu/dnsx/model/dao"
)

func TypeNS(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	for _, rr := range rrs {
		msg.Answer = append(msg.Answer, &dns.NS{
			Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			Ns:  rr.Value,
		})
	}
	return nil
}
