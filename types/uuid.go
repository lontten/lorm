package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"strings"
)

type UUID struct {
	uuid.UUID
}

func (u UUID) MarshalJSON() ([]byte, error) {
	all := strings.ReplaceAll(u.String(), "-", "")
	rs := []byte(fmt.Sprintf(`"%s"`, all))

	return rs, nil

}

func (u *UUID) UnmarshalJSON(src []byte) error {
	if len(src) != 34 {
		return errors.Errorf("invalid length for UUID: %v", len(src))
	}
	fromString, err := uuid.FromString(string(src[1 : len(src)-1]))
	if err != nil {
		return err
	}
	*u = UUID{fromString}
	return err
}

// Value insert timestamp into mysql need this function.

func (u UUID) Value() (driver.Value, error) {
	return u.UUID.String(), nil
}

// Scan valueof time.Time
func (u *UUID) Scan(v interface{}) error {
	value, ok := v.(string)
	if ok {
		*u = UUID{uuid.FromStringOrNil(value)}
		return nil
	}
	return fmt.Errorf("can not convert %v to uuid", v)
}

func String2UUID(v string) UUID {
	return UUID{uuid.FromStringOrNil(v)}
}

type UUIDList []UUID

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p UUIDList) Value() (driver.Value, error) {
	var k []UUID
	k = p
	marshal, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	var s = string(marshal)
	s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
	return s, nil
}

// Scan 实现方法
func (p *UUIDList) Scan(data interface{}) error {
	array := pgtype.UUIDArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []UUID
	list = make([]UUID, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = UUID{element.Bytes}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}