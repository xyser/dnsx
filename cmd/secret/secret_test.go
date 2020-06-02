package secret

import (
	"fmt"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestB(t *testing.T) {
	dnskey257 := &dns.DNSKEY{
		Hdr: dns.RR_Header{
			Name:   "cloudflare.com.",
			Rrtype: dns.TypeDNSKEY,
			Class:  dns.ClassINET,
			Ttl:    3115,
		},
		Flags:     257,
		Protocol:  3,
		Algorithm: dns.ECDSAP256SHA256,
		PublicKey: "mdsswUyr3DPW132mOi8V9xESWE8jTo0dxCjjnopKl+GqJxpVXckHAeF+KkxLbxILfDLUT0rAK9iUzy1L53eKGQ==",
	}

	dnskey256 := &dns.DNSKEY{
		Hdr: dns.RR_Header{
			Name:   "cloudflare.com.",
			Rrtype: dns.TypeDNSKEY,
			Class:  dns.ClassINET,
			Ttl:    3115,
		},
		Flags:     256,
		Protocol:  3,
		Algorithm: dns.ECDSAP256SHA256,
		PublicKey: "oJMRESz5E4gYzS/q6XDrvU1qMPYIjCWzJaOau8XNEZeqCYKD5ar0IRd8KqXXFJkqmVfRvMGPmM1x8fGAa2XhSA==",
	}

	fmt.Println(dnskey257)
	fmt.Println(dnskey256)
	fmt.Println(dnskey256.KeyTag())

	olt, _ := time.ParseInLocation("20060102150405", "20200623071422", time.UTC)
	olt2, _ := time.ParseInLocation("20060102150405", "20200424071422", time.UTC)
	sig := &dns.RRSIG{
		Hdr: dns.RR_Header{
			Name:   "cloudflare.com.",
			Rrtype: dns.TypeRRSIG,
			Class:  dns.ClassINET,
			Ttl:    3115,
		},
		TypeCovered: dns.TypeDNSKEY,
		Algorithm:   dns.ECDSAP256SHA256,
		Labels:      2,
		OrigTtl:     3600,
		Expiration:  uint32(olt.Unix()),
		Inception:   uint32(olt2.Unix()),
		KeyTag:      dnskey257.KeyTag(),
		SignerName:  dnskey257.Hdr.Name,
		Signature:   "M4Gm5gntEncnqRht2ALaUd8tu/NHqIWTCIlu/anbncHYYE3qJx00sVTFdkuxpaRRQRRI7HXd/dJjxZi/2FWLJg==",
	}

	fmt.Println(sig.Verify(dnskey257, []dns.RR{dnskey256}))
	fmt.Println(sig)
}

func TestVerteiltesysteme_dot_net_rrsig(t *testing.T) {
	key257, _ := dns.NewRR("cloudflare.com.\t\t3599\tIN\tDNSKEY\t257 3 13 mdsswUyr3DPW132mOi8V9xESWE8jTo0dxCjjnopKl+GqJxpVXckHAeF+ KkxLbxILfDLUT0rAK9iUzy1L53eKGQ==")
	key256, _ := dns.NewRR("cloudflare.com.\t\t3599\tIN\tDNSKEY\t256 3 13 oJMRESz5E4gYzS/q6XDrvU1qMPYIjCWzJaOau8XNEZeqCYKD5ar0IRd8 KqXXFJkqmVfRvMGPmM1x8fGAa2XhSA==")
	sigs, _ := dns.NewRR("cloudflare.com.\t\t3599\tIN\tRRSIG\tDNSKEY 13 2 3600 20200623071422 20200424071422 2371 cloudflare.com. M4Gm5gntEncnqRht2ALaUd8tu/NHqIWTCIlu/anbncHYYE3qJx00sVTF dkuxpaRRQRRI7HXd/dJjxZi/2FWLJg==")

	zsk := key257.(*dns.DNSKEY)
	rsk := key256.(*dns.DNSKEY)

	sin := sigs.(*dns.RRSIG)

	t.Log(rsk)
	t.Log(sin)
	t.Log(zsk)
	if e := sin.Verify(zsk, []dns.RR{rsk}); e != nil {
		t.Errorf("cannot verify RRSIG with keytag [%d]. Cause [%s]", rsk.KeyTag(), e.Error())
	}
}

func TestAAA(t *testing.T) {

	//key257, _ := dns.NewRR("cloudflare.com.\t\t3599\tIN\tDNSKEY\t257 3 13 mdsswUyr3DPW132mOi8V9xESWE8jTo0dxCjjnopKl+GqJxpVXckHAeF+ KkxLbxILfDLUT0rAK9iUzy1L53eKGQ==")
	key256, _ := dns.NewRR("cloudflare.com.\t\t3599\tIN\tDNSKEY\t256 3 13 oJMRESz5E4gYzS/q6XDrvU1qMPYIjCWzJaOau8XNEZeqCYKD5ar0IRd8 KqXXFJkqmVfRvMGPmM1x8fGAa2XhSA==")
	sigs, _ := dns.NewRR("cloudflare.com.\t\t108\tIN\tRRSIG\tA 13 2 300 20200602115010 20200531095010 34505 cloudflare.com. 1ntvMbesPT/IeNKAzB+L2UZ0OCvAqbszfQeu4Ijn0rQJVgZfs4fdKEvu Gv6a0fZxrMbAeU93vTXDMLaVeHVg7w==")
	a1, _ := dns.NewRR("cloudflare.com.\t\t108\tIN\tA\t104.17.175.85")
	a2, _ := dns.NewRR("cloudflare.com.\t\t108\tIN\tA\t104.17.176.85")

	//zsk := key257.(*dns.DNSKEY)
	rsk := key256.(*dns.DNSKEY)

	sin := sigs.(*dns.RRSIG)
	a1s := a1.(*dns.A)
	a2s := a2.(*dns.A)

	t.Log(rsk)
	t.Log(sin)
	t.Log(a1s)
	t.Log(a2s)

	if e := sin.Verify(rsk, []dns.RR{a1s, a2s}); e != nil {
		t.Errorf("cannot verify RRSIG with keytag [%d]. Cause [%s]", rsk.KeyTag(), e.Error())
	}

}
