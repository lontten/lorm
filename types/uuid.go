package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
)

type UUIDList []uuid.UUID

// Value 实现方法
func (p UUIDList) Value() (driver.Value, error) {
	var k []uuid.UUID
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
func (p *UUIDList) Scan(data interface{}) error {
	array := pgtype.UUIDArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []uuid.UUID
	list = make([]uuid.UUID, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = element.Bytes
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}
