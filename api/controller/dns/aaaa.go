package dns

import (
	"net"

	"github.com/miekg/dns"

	"github.com/xyser/dnsx/model/dao"
)

// TypeAAAA query aaaa
func TypeAAAA(msg *dns.Msg) error {
	for _, question := range msg.Question {
		rrs, err := dao.GetRecordByNameAndType(question.Name, question.Qtype)
		if err != nil {
			return err
		}
		for _, rr := range rrs {
			msg.Answer = append(msg.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: question.Name, Rrtype: question.Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
				AAAA: net.ParseIP(rr.Value),
			})
		}
	}

	return nil
}
