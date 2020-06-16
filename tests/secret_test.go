package tests

import (
	"testing"

	"github.com/miekg/dns"
)

func TestVerificationCloudflareDotComDNSKEYRRSIG(t *testing.T) {
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
		//t.Errorf("cannot verify RRSIG with keytag [%d]. Cause [%s]", rsk.KeyTag(), e.Error())
	}
}

func TestVerificationCloudflareDotComARRSIG(t *testing.T) {
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
