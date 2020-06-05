package dns

import (
	"strings"

	"github.com/miekg/dns"

	"github.com/dingdayu/dnsx/model/dao"
)

// TypeCAA query caa
func TypeCAA(msg *dns.Msg) error {
	domain := msg.Question[0].Name
	rrs, err := dao.GetRecordByNameAndType(domain, msg.Question[0].Qtype)
	if err != nil {
		return err
	}
	// CAA 协议解释: https://support.dnsimple.com/articles/caa-record/
	// RFC 文档: https://tools.ietf.org/html/rfc6844#section-3
	for _, rr := range rrs {
		// 标志位[0,1] 标签位[issue,issuewild,iodef] 机构域
		// domain.com. CAA 0 iodef mailto:admin@domain.com
		vs := strings.Split(rr.Value, " ")
		// 不够3位 不返回
		if len(vs) < 3 {
			continue
		}

		flag := 0
		if vs[0] == "1" {
			flag = 1
		}
		msg.Answer = append(msg.Answer, &dns.CAA{
			Hdr:  dns.RR_Header{Name: msg.Question[0].Name, Rrtype: msg.Question[0].Qtype, Class: dns.ClassINET, Ttl: rr.TTL},
			Flag: uint8(flag), // 标志位，严格校验 Tag 标签位
			// Tag 标签位
			// issue: 显式地授权单个证书颁发机构为主机名颁发证书（任何类型）。
			// issuewild: 显式地授权单个证书颁发机构为主机名颁发通配符证书（并且仅通配符）。
			// iodef: 指定证书颁发机构可以向其报告策略违规的URL。使用了事件对象描述交换格式（IODEF）格式
			Tag:   vs[1], // 标签位，
			Value: vs[2],
		})
	}
	return nil
}
