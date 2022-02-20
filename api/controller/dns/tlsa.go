package dns

import (
	"strconv"
	"strings"

	"github.com/miekg/dns"

	"github.com/xyser/dnsx/model/dao"
)

// TypeTLSA TLSA记录格式: 保存证书关联数据
func TypeTLSA(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	// TLSA 协议解释: https://www.dynu.com/Resources/DNS-Records/TLSA-Record
	for _, rr := range rrs {
		vs := strings.Split(rr.Value, " ")
		// 不够3位 不返回
		if len(vs) < 4 {
			continue
		}
		usage, _ := strconv.ParseUint(vs[0], 8, 8)
		selector, _ := strconv.ParseUint(vs[1], 8, 8)
		match, _ := strconv.ParseUint(vs[2], 8, 8)

		msg.Answer = append(msg.Answer, &dns.TLSA{
			Hdr:          dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			Usage:        uint8(usage),    // 证书使用情况
			Selector:     uint8(selector), // 选择器,0-完整证书,1-使用主题公钥
			MatchingType: uint8(match),    // 匹配类型, 0-无哈希, 1-所选内容的SHA-256哈希,2-所选内容的SHA-512哈希
			Certificate:  "",              // 证书关联数据
		})
	}
	return nil
}
