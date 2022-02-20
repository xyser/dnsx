package dns

import (
	"strings"

	"github.com/miekg/dns"

	"github.com/xyser/dnsx/model/dao"
)

// TypeHINFO hardware type and operating system (OS) information
// testhinfo 90 IN HINFO "INTEL-386" "Windows"
func TypeHINFO(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	for _, rr := range rrs {
		vs := strings.Split(rr.Value, " ")
		// 不够3位 不返回
		if len(vs) < 2 {
			continue
		}
		msg.Answer = append(msg.Answer, &dns.HINFO{
			Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			Cpu: vs[0],
			Os:  vs[1],
		})
	}
	return nil
}
