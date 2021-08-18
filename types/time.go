package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

type Time time.Time

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
	var s = ""
	switch v := v.(type) {
	case string:
		s = v[:8]
	case []byte:
		s = string(v)[:8]
	case time.Time:
		*t = Time(v)
		return nil
	default:
		return fmt.Errorf("can not convert %v to timestamp", v)
	}
	now, err := time.Parse(`15:04:05`, s)
	if err != nil {
		return err
	}
	*t = Time(now)
	return nil
}

type TimeList []Time

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p TimeList) Value() (driver.Value, error) {
	var k []Time
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
func (p *TimeList) Scan(data interface{}) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []Time
	list = make([]Time, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = Time(element.Time)
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

//date
type Date time.Time

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
	return t, nil
}

// Scan valueof jstime.Time
func (t *Date) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Date(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
func (t Date) ToGoTime() time.Time {
	return time.Unix(time.Time(t).Unix(), 0)
}

type DateList []Date

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p DateList) Value() (driver.Value, error) {
	var k []Date
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
func (p *DateList) Scan(data interface{}) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []Date
	list = make([]Date, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = Date(element.Time)
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

//datetime
type DateTime time.Time

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
	return t, nil
}

// Scan valueof jstime.Time
func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t DateTime) ToGoTime() time.Time {
	return time.Unix(time.Time(t).Unix(), 0)
}

type DateTimeList []DateTime

// gorm 自定义结构需要实现 Value Scan 两个方法
// Value 实现方法
func (p DateTimeList) Value() (driver.Value, error) {
	var k []DateTime
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
func (p *DateTimeList) Scan(data interface{}) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []DateTime
	list = make([]DateTime, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = DateTime(element.Time)
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

func (d Date) AddTime(t Time) DateTime {
	d2 := time.Time(d)
	t2 := time.Time(t)

	return DateTime(time.Date(
		d2.Year(),
		d2.Month(),
		d2.Day(),
		t2.Hour(),
		t2.Minute(),
		t2.Second(), 0, nil,
	))
}

func (t Time) AddData(d Date) DateTime {
	d2 := time.Time(d)
	t2 := time.Time(t)

	return DateTime(time.Date(
		d2.Year(),
		d2.Month(),
		d2.Day(),
		t2.Hour(),
		t2.Minute(),
		t2.Second(), 0, nil,
	))
}

//datetime
type AutoDateTime struct {
	time.Time
}

func (t AutoDateTime) MarshalJSON() ([]byte, error) {
	var tune string
	if t.Year() == 0 && t.Month() == time.January && t.Day() == 1 {
		tune = t.Format(`"15:04:05"`)
	} else {
		tune = t.Format(`"2006-01-02"`)
	}
	return []byte(tune), nil
}

func (t *AutoDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = AutoDateTime{Time: now}
	return err
}

// Value insert timestamp into mysql need this function.
func (t AutoDateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof jstime.Time
func (t *AutoDateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = AutoDateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
