package entity

import (
	"github.com/miekg/dns"
)

// JsonMsgQuestion json question
type JsonMsgQuestion struct {
	Name string `json:"name"` // FQDN with trailing dot
	Type uint16 `json:"type"` // SPF - Standard DNS RR type
}

// JsonMsgAnswer json answer
type JsonMsgAnswer struct {
	Name string `json:"name"` // Always matches name in Question
	Type uint16 `json:"type"` // SPF - Standard DNS RR type
	TTL  uint32 `json:"TTL"`  // Record's time-to-live in seconds
	Data string `json:"data"` // Data for SPF - quoted string
}

// JsonMsgAdditional json additional
type JsonMsgAdditional struct {
}

// JsonMsgResponse json response
type JsonMsgResponse struct {
	Status           int  // NOERROR - Standard DNS response code (32 bit integer).
	TC               bool // Whether the response is truncated
	RD               bool // Always true for Google Public DNS
	RA               bool // Always true for Google Public DNS
	AD               bool // Whether all response data was validated with DNSSEC
	CD               bool // Whether the client asked to disable DNSSEC
	Question         []JsonMsgQuestion
	Answer           []JsonMsgAnswer
	Additional       []JsonMsgAdditional
	EDNSClientSubnet string `json:"edns_client_subnet"`
	Comment          string `json:"Comment"` // Uncached responses are attributed to the authoritative name server
}

// MsgToJson msg to json message
func MsgToJson(msg *dns.Msg) (resp JsonMsgResponse) {
	resp.Status = msg.Rcode
	resp.TC = msg.Truncated
	resp.RD = msg.RecursionDesired
	resp.RA = msg.RecursionAvailable
	resp.AD = msg.AuthenticatedData
	resp.CD = msg.CheckingDisabled
	for _, v := range msg.Question {
		resp.Question = append(resp.Question, JsonMsgQuestion{
			Name: v.Name,
			Type: v.Qtype,
		})
	}
	for _, v := range msg.Answer {
		resp.Answer = append(resp.Answer, JsonMsgAnswer{
			Name: v.Header().Name,
			Type: v.Header().Rrtype,
			TTL:  v.Header().Ttl,
			Data: v.Header().String(),
		})
	}
	return resp
}
