package model

// 4.3.2.6.in-addr.arpa.
// 167.10.168.192.IN-ADDR.ARPA.

var DomainsToAddresses = map[string]string{
	"google.com.": "1.2.3.4",
	"ddy.":        "104.198.14.52",
}

type Domain struct {
}

// (Resource Record) http://dns-learning.twnic.net.tw/bind/intro6.html
type RR struct {
	Type   int
	Name   string
	Value  string
	Filter string
	TTL    int
}

func Add() {

}

func GetDomain(domain string) string {
	return DomainsToAddresses[domain]
}

func GetAll() map[string]string {
	return DomainsToAddresses
}
