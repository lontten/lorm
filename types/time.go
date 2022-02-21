package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time time.Time

func NowTime() Time {
	return Time(time.Now())
}

func NowTimeP() *Time {
	now := time.Now()
	return (*Time)(&now)
}

func (t Time) MarshalJSON() ([]byte, error) {

	tune := time.Time(t).Format(`"15:04:05"`)
	return []byte(tune), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"15:04:05"`, string(data), time.Local)
	*t = Time(now)
	return err
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t, nil
}

func (t *Time) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	var s = ""
	switch v := v.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	case time.Time:
		*t = Time(v)
	case Time:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to Time", v)
	}
	if len(s) < 8 {
		return nil
	}
	now, err := time.Parse(`15:04:05`, s[:8])
	if err != nil {
		return err
	}
	*t = Time(now)
	return nil
}

//date
type Date time.Time

func NowDate() Date {
	return Date(time.Now())
}

func NowDateP() *Date {
	now := time.Now()
	return (*Date)(&now)
}

func (t Date) MarshalJSON() ([]byte, error) {
	tune := time.Time(t).Format(`"2006-01-02"`)
	return []byte(tune), nil
}

func (t *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*t = Date(now)
	return err
}

// Value insert timestamp into mysql need this function.
func (t Date) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(t), nil
}

// Scan valueof jstime.Time
func (t *Date) Scan(v interface{}) error {

	value, ok := v.(time.Time)
	if ok {
		*t = Date(value)
		return nil
	}

	value2, ok2 := v.(Date)
	if ok2 {
		*t = value2
		return nil
	}

	return fmt.Errorf("can not convert %v to types.Date", v)

}
func (t Date) ToGoTime() time.Time {
	return time.Unix(time.Time(t).Unix(), 0)
}

//datetime
type DateTime time.Time

func NowDateTime() DateTime {
	return DateTime(time.Now())
}

func NowDateTimeP() *DateTime {
	now := time.Now()
	return (*DateTime)(&now)
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	tune := time.Time(t).Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = DateTime(now)
	return err
}

// Value insert timestamp into mysql need this function.
func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(t), nil
}

// Scan valueof jstime.Time
func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime(value)
		return nil
	}

	value2, ok2 := v.(DateTime)
	if ok2 {
		*t = value2
		return nil
	}

	return fmt.Errorf("can not convert %v to types.DateTime", v)

}

func (t DateTime) ToGoTime() time.Time {
	return time.Unix(time.Time(t).Unix(), 0)
}

//datetime
type AutoDateTime time.Time

func (t AutoDateTime) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	var tune string
	if tt.Year() == 0 && tt.Month() == time.January && tt.Day() == 1 {
		tune = tt.Format(`"15:04:05"`)
	} else {
		tune = tt.Format(`"2006-01-02"`)
	}
	return []byte(tune), nil
}

func (t *AutoDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = AutoDateTime(now)
	return err
}

// Value insert timestamp into mysql need this function.
func (t AutoDateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(t), nil
}

// Scan valueof jstime.Time
func (t *AutoDateTime) Scan(v interface{}) error {
	var s = ""
	switch v := v.(type) {
	case string:
		s = v[:8]
	case []byte:
		s = string(v)[:8]
	case time.Time:
		*t = AutoDateTime(v)
	case Time:
		*t = AutoDateTime(v)
	case Date:
		*t = AutoDateTime(v)
	case AutoDateTime:
		*t = v
		return nil
	default:
		return fmt.Errorf("can not convert %v to types.AutoDateTime", v)
	}
	now, err := time.Parse(`2006-01-02 15:04:05`, s)
	if err != nil {
		return err
	}
	*t = AutoDateTime(now)
	return nil
}
