package dns

import (
	"strconv"
	"strings"

	"github.com/miekg/dns"

	"github.com/dingdayu/dnsx/model/dao"
)

// TypeURI provide a means to resolve hostnames to URIs that can be used by various applications.
// For resolution the application needs to know both the hostname and the protocol that the URI is to be used for.
//
// It is an alternative to the SRV record. Similar to the SRV record,
// it returns both weight and priority values which can be used to select an appropriate URI from multiple results.
// However, unlike the SRV no port number is returned in the URI record since this information
// is contained (if applicable) in the URI string.
//
// So the returned URI strings from a URI record query can be used directly by the requested application,
// while the usable URI has to be assembled from the search and results information of an SRV record query.
//
// _ftp._tcp.example.com. 3600 IN URI 10 1 "ftp://ftp.example.com/public"
func TypeURI(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	for _, rr := range rrs {
		// 优先级 权重 端口 目标地址 1 1 80 www.baidu.com
		vs := strings.Split(rr.Value, " ")
		var priority, weight uint64
		if len(vs) == 3 {
			priority, _ = strconv.ParseUint(vs[0], 16, 16)
			weight, _ = strconv.ParseUint(vs[1], 16, 16)
		} else if len(vs) == 2 {
			priority = uint64(rr.Priority)
			weight, _ = strconv.ParseUint(vs[0], 16, 16)
		}

		// URI 协议解释: https://www.dynu.com/Resources/DNS-Records/URI-Record
		msg.Answer = append(msg.Answer, &dns.URI{
			Hdr:      dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			Priority: uint16(priority), // 优先级
			Weight:   uint16(weight),   // 权重
			Target:   vs[len(vs)-1],    // 对应目标地址,可以是域名或IP
		})
	}
	return nil
}
