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

// NullUUID represents a uuid.UUID that may be null.
// NullUUID implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUUID struct {
	UUID  uuid.UUID
	Valid bool // Valid is true if UUID is not NULL
}

// Scan implements the Scanner interface.
func (u *NullUUID) Scan(v interface{}) error {
	if v == nil {
		u.UUID, u.Valid = uuid.UUID{}, false
		return nil
	}
	u.Valid = true
	value, ok := v.(string)
	if ok {
		u.UUID = uuid.FromStringOrNil(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Value implements the driver Valuer interface.
func (u NullUUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.UUID, nil
}

func (u NullUUID) IsNull() bool {
	return !u.Valid
}

func (u NullUUID) MarshalJSON() ([]byte, error) {
	var str = ""
	if !u.Valid {
		return []byte(fmt.Sprintf(`"%s"`, str)), nil
	}
	str = strings.ReplaceAll(u.UUID.String(), "-", "")
	return []byte(fmt.Sprintf(`"%s"`, str)), nil

}

func (u *NullUUID) UnmarshalJSON(src []byte) error {
	if len(src) != 34 {
		return errors.Errorf("invalid length for UUID: %v", len(src))
	}
	orNil := uuid.FromStringOrNil(string(src[1 : len(src)-1]))
	u.UUID, u.Valid = orNil, orNil != uuid.Nil
	return nil
}

func (u *NullUUID) Set(v string) {
	orNil := uuid.FromStringOrNil(v)
	u.UUID, u.Valid = orNil, orNil != uuid.Nil
}

func (u *NullUUID) SetUUID(v uuid.UUID) {
	u.UUID, u.Valid = v, v != uuid.Nil
}

func String2NullUUID(v string) NullUUID {
	var u NullUUID
	orNil := uuid.FromStringOrNil(v)
	u.UUID, u.Valid = orNil, orNil != uuid.Nil
	return u
}

func String2UUID(v string) UUID {
	return UUID{uuid.FromStringOrNil(v)}
}

func FromUUID(uuid UUID) NullUUID {
	return NullUUID{UUID: uuid.UUID, Valid: true}
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
