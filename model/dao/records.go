package dao

import (
	"errors"

	"github.com/miekg/dns"

	"github.com/dingdayu/dnsx/internal/asset"
	"github.com/dingdayu/dnsx/internal/mysql"
)

// Record struct
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

// TypeEnum record type enum
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
	// ErrTypeEnumKey type enum not find error
	ErrTypeEnumKey = errors.New("type enum key error")
)

// TableName db table name
func (Record) TableName() string {
	return "records"
}

// InitRecord 检查表是否存在并创建表
func InitRecord() {
	db := mysql.GetDB()
	if !db.HasTable(Record{}) {
		if sql, err := asset.Asset("scripts/sql/record.sql"); err == nil {
			db.Exec(string(sql))
		}
	}
}

// CreateRecord insert db
func CreateRecord(rr *Record) error {
	return mysql.GetDB().Create(&rr).Error
}

// GetNameRecord query name by db
func GetNameRecord(name string) (rrs []Record, err error) {
	err = mysql.GetDB().Where("name = ?", name).Find(&rrs).Error
	return
}

// GetRecord query db by where
func GetRecord(where interface{}) (rrs []Record, err error) {
	err = mysql.GetDB().Where(where).Find(&rrs).Error
	return rrs, err
}

// GetRecordByNameAndType query db by name and type
func GetRecordByNameAndType(name string, qtype uint16) (rrs []Record, err error) {
	if types, ok := TypeEnum[qtype]; !ok {
		return rrs, ErrTypeEnumKey
	} else {
		err = mysql.GetDB().Where("name = ? AND type = ?", name, types).Find(&rrs).Error
	}
	return
}
