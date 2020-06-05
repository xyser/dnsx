package entity

// RR (Resource Record) http://dns-learning.twnic.net.tw/bind/intro6.html
type RR struct {
	Type   int
	Name   string
	Value  string
	Filter string
	TTL    int
}
