package dao

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateFormat date time
const DateFormat = "2006-01-02"

// LocalDate local date
type LocalDate struct {
	time.Time
}

// MarshalJSON LocalDate 序列号
func (t LocalDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Format(DateFormat))), nil
}

// Value LocalDate 转 time
func (t LocalDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan Gorm 扫描时的数据赋值
func (t *LocalDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
