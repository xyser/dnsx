package dns

import (
	"github.com/miekg/dns"

	"github.com/dingdayu/dnsx/model/dao"
	"github.com/dingdayu/dnsx/pkg/network"
)

// TypePTR query ptr
func TypePTR(msg *dns.Msg) error {
	msg.Authoritative = true
	ip := network.PTRToIP([]byte(msg.Question[0].Name))

	rrs, err := dao.GetRecord(map[string]interface{}{"type": "a", "value": ip})
	if err != nil {
		return err
	}
	for _, rr := range rrs {
		msg.Answer = append(msg.Answer, &dns.PTR{
			Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			Ptr: rr.Name,
		})
	}
	return nil
}
