package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type DecimalList []decimal.Decimal

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p DecimalList) Value() (driver.Value, error) {
	var k []decimal.Decimal
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
func (p *DecimalList) Scan(data any) error {
	array := pgtype.VarcharArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []decimal.Decimal
	list = make([]decimal.Decimal, len(array.Elements))
	for i, element := range array.Elements {
		fromString, err := decimal.NewFromString(element.String)
		if err != nil {
			return err
		}
		list[i] = fromString
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
