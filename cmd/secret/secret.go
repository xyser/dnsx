package secret

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

// Cmd cmd example
var Cmd = &cobra.Command{
	Use:   "secret",
	Short: "示例",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := new(dns.DNSKEY)
		key.Hdr = dns.RR_Header{
			Name:   "dingxiaoyu.com.",
			Rrtype: dns.TypeDNSKEY,
			Class:  dns.ClassINET,
			Ttl:    14400,
		}
		key.Flags = 256
		key.Protocol = 3
		key.Algorithm = dns.ECDSAP256SHA256

		//privkey,_ := key.Generate(256)
		//prive := key.PrivateKeyString(privkey)

		prive := "Private-key-format: v1.3\nAlgorithm: 13 (ECDSAP256SHA256)\nPrivateKey: D8GGLJzXo/O3H16PnpGrBNZF18a+w0MTr1Mvnm/25wc=\n"

		privkey, _ := key.NewPrivateKey(prive)

		// The record we want to sign
		srv := new(dns.A)
		srv.Hdr = dns.RR_Header{Name: "qq.dingxiaoyu.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 14400}
		srv.A = net.ParseIP("127.0.0.1")

		// 产生 DS 记录
		//ds := key.ToDS(dns.SHA256).String()
		//fmt.Println(ds )

		sig := new(dns.RRSIG)
		sig.Hdr = dns.RR_Header{Name: "dingxiaoyu.com.", Rrtype: dns.TypeRRSIG, Class: dns.ClassINET, Ttl: 14400}
		sig.TypeCovered = srv.Hdr.Rrtype
		sig.Labels = uint8(dns.CountLabel(srv.Hdr.Name)) // works for all 3
		sig.OrigTtl = srv.Hdr.Ttl
		sig.Expiration = 1296534305 // date -u '+%s' -d"2011-02-01 04:25:05"
		sig.Inception = 1293942305  // date -u '+%s' -d"2011-01-02 04:25:05"
		sig.KeyTag = key.KeyTag()   // Get the keyfrom the Key
		sig.SignerName = key.Hdr.Name
		sig.Algorithm = dns.ECDSAP256SHA256

		fmt.Println(key.String())
		fmt.Println(key.PrivateKeyString(privkey))
		fmt.Println(privkey)

	},
}
