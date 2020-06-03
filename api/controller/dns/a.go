package dns

import (
	"net"

	"github.com/miekg/dns"

	"github.com/dingdayu/dnsx/model/dao"
)

func TypeA(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	for _, rr := range rrs {
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			A:   net.ParseIP(rr.Value),
		})
	}
	return nil
}
