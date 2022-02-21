package types

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type UUID uuid.UUID

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u UUID) MarshalJSON() ([]byte, error) {
	all := strings.ReplaceAll(uuid.UUID(u).String(), "-", "")
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
	*u = UUID(fromString)
	return err
}

// Value insert timestamp into mysql need this function.

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).String(), nil
}

// Scan valueof time.Time
func (u *UUID) Scan(v interface{}) error {
	value, ok := v.(string)
	if ok {
		*u = UUID(uuid.FromStringOrNil(value))
		return nil
	}
	return fmt.Errorf("can not convert %v to uuid", v)
}

func Str2UUIDMust(v string) UUID {
	return UUID(uuid.FromStringOrNil(v))
}

func Str2UUIDMustP(v string) *UUID {
	fromString, err := uuid.FromString(v)
	if err != nil {
		return nil
	}
	u := UUID(fromString)
	return &u
}

func NewV4() UUID {
	return UUID(uuid.NewV4())
}

func NewV4P() *UUID {
	u := UUID(uuid.NewV4())
	return &u
}

func Str2UUID(v string) (UUID, error) {
	id, err := uuid.FromString(v)
	if err != nil {
		return UUID{}, err
	}
	return UUID(id), nil
}
