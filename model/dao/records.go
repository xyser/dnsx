package dao

import (
	"errors"
	"github.com/miekg/dns"
)

type Record struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Value     string     `json:"value"`
	TTL       uint32     `json:"ttl"`
	Priority  int        `json:"priority"`
	CreatedAt *LocalTime `json:"created_at"`
	UpdatedAt *LocalTime `json:"updated_at"`
}

var TypeEnum = map[uint16]string{
	dns.TypeA:     "a",
	dns.TypeAAAA:  "aaaa",
	dns.TypeCNAME: "cname",
	dns.TypePTR:   "ptr",
	dns.TypeNS:    "ns",
	dns.TypeSRV:   "srv",
	dns.TypeTXT:   "txt",
	dns.TypeCAA:   "caa",
	dns.TypeHINFO: "hinfo",
	dns.TypeTLSA:  "tlsa",
	dns.TypeURI:   "uri",
}

var (
	ErrTypeEnumKey = errors.New("type enum key error")
)

// TableName db table name
func (Record) TableName() string {
	return "records"
}

// CreateRecord insert db
func CreateRecord(rr *Record) error {
	return db.Create(&rr).Error
}

// GetNameRecord query name by db
func GetNameRecord(name string) (rrs []Record, err error) {
	err = db.Where("name = ?", name).Find(&rrs).Error
	return
}

// GetRecord query db by where
func GetRecord(where interface{}) (rrs []Record, err error) {
	err = db.Where(where).Find(&rrs).Error
	return
}

// GetRecordByNameAndType query db by name and type
func GetRecordByNameAndType(name string, qtype uint16) (rrs []Record, err error) {
	if types, ok := TypeEnum[qtype]; !ok {
		return rrs, ErrTypeEnumKey
	} else {
		err = db.Where("name = ? AND type = ?", name, types).Find(&rrs).Error
	}
	return
}
