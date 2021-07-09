package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/jackc/pgtype"
)

type StringList []string

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p StringList) Value() (driver.Value, error) {
	var k []string
	k = p
	marshal, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	var s = string(marshal)
	if s != "null" {
		s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
	}
	return s, nil
}

func (p StringList) IsNull() bool {
	return len(p) == 0
}

// Scan 实现方法
func (p *StringList) Scan(data interface{}) error {
	array := pgtype.VarcharArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []string
	list = make([]string, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = element.String
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
