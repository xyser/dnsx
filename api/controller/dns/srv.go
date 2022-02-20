package dns

import (
	"strconv"
	"strings"

	"github.com/miekg/dns"

	"github.com/xyser/dnsx/model/dao"
)

// TypeSRV is a specification of servers in the Domain Name System by hostname and port number.
// With an SRV record, it is possible to make a server discoverable and designate high priority
// and high availability servers using a single domain without having to know the exact address
// of the servers.
func TypeSRV(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	for _, rr := range rrs {
		// 优先级 权重 端口 目标地址 1 1 80 www.baidu.com
		vs := strings.Split(rr.Value, " ")
		var priority, weight, port uint64
		if len(vs) == 4 {
			priority, _ = strconv.ParseUint(vs[0], 16, 16)
			weight, _ = strconv.ParseUint(vs[1], 16, 16)
			port, _ = strconv.ParseUint(vs[2], 16, 16)
		} else if len(vs) == 3 {
			priority = uint64(rr.Priority)
			weight, _ = strconv.ParseUint(vs[0], 16, 16)
			port, _ = strconv.ParseUint(vs[1], 16, 16)
		}
		msg.Answer = append(msg.Answer, &dns.SRV{
			Hdr:      dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			Priority: uint16(priority), // 优先级
			Weight:   uint16(weight),   // 权重
			Port:     uint16(port),     // 端口
			Target:   vs[len(vs)-1],    // 对应目标地址,可以是域名或IP
		})
	}
	return nil
}
