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
	} else {
		s = "{}"
	}
	return s, nil
}

// Scan 实现方法
func (p *StringList) Scan(data any) error {
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

func (p StringList) Len() int {
	return len(p)
}

// 实现Less方法
func (p StringList) Less(i, j int) bool {
	return p[i] < p[j]
}

// 实现Swap方法
func (p StringList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
