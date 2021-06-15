package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/jackc/pgtype"
)

type NullBool sql.NullBool

// Scan implements the Scanner interface.
func (n *NullBool) Scan(value interface{}) error {
	return (*sql.NullBool)(n).Scan(value)
}

// Value implements the driver Valuer interface.

func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

func (n NullBool) IsNull() bool {
	return !n.Valid
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

func (n *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Bool)
	if err == nil {
		n.Valid = true
	}
	return err
}

func (b *NullBool) Set(v bool) {
	b.Bool, b.Valid = v, true
}

func From(b bool) NullBool {
	return NullBool{
		Bool:  b,
		Valid: true,
	}
}

func Pg2Arr(v pgtype.BoolArray) []bool {
	var arr []bool
	for _, element := range v.Elements {
		if element.Status == pgtype.Present {
			arr = append(arr, element.Bool)
		}
	}
	return arr
}

func Arr2Pg(arr []bool) pgtype.BoolArray {
	var list = pgtype.BoolArray{}
	list.Set(arr)
	return list
}
