package dns

import (
	"net"

	"github.com/miekg/dns"

	"github.com/xyser/dnsx/model/dao"
)

// TypeAAAA query aaaa
func TypeAAAA(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	for _, rr := range rrs {
		msg.Answer = append(msg.Answer, &dns.AAAA{
			Hdr:  dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			AAAA: net.ParseIP(rr.Value),
		})
	}
	return nil
}
