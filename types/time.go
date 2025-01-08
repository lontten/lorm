package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

type Time struct {
	time.Time
}

func NowTime() Time {
	return Time{time.Now()}
}

func NowTimeP() *Time {
	return &Time{time.Now()}
}
func TimeOf(t time.Time) Time {
	timeOnly := time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return Time{timeOnly}
}
func TimePOf(t time.Time) *Time {
	timeOnly := time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return &Time{timeOnly}
}
func (t Time) ToGoTime() time.Time {
	return time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}
func (t Time) ToDateTime() DateTime {
	return DateTime{t.ToGoTime()}
}
func (t Time) ToDateTimeP() *DateTime {
	return &DateTime{t.ToGoTime()}
}
func (t Time) ToString() string {
	return t.Time.Format(`15:04:05`)
}

func (t Time) MarshalJSON() ([]byte, error) {
	tune := t.Time.Format(`"15:04:05"`)
	return []byte(tune), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"15:04:05"`, string(data), time.Local)
	*t = Time{now}
	return err
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *Time) Scan(v any) error {
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
		*t = Time{v}
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
	*t = Time{Time: now}
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
func (p *TimeList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []Time
	list = make([]Time, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = Time{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

// date
type Date struct {
	time.Time
}

func NowDate() Date {
	return Date{time.Now()}
}

func NowDateP() *Date {
	return &Date{time.Now()}
}
func DateOf(t time.Time) Date {
	dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return Date{dateOnly}
}
func DatePOf(t time.Time) *Date {
	dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &Date{dateOnly}
}

func (t Date) ToGoTime() time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
func (t Date) ToDateTime() DateTime {
	return DateTime{t.ToGoTime()}
}
func (t Date) ToDateTimeP() *DateTime {
	return &DateTime{t.ToGoTime()}
}

func (t Date) ToString() string {
	return t.Format(`2006-01-02`)
}

func (t Date) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02"`)
	return []byte(tune), nil
}

func (t *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*t = Date{Time: now}
	return err
}

// Value insert timestamp into mysql need this function.
func (t Date) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof jstime.Time
func (t *Date) Scan(v any) error {

	value, ok := v.(time.Time)
	if ok {
		*t = Date{value}
		return nil
	}

	value2, ok2 := v.(Date)
	if ok2 {
		*t = value2
		return nil
	}

	return fmt.Errorf("can not convert %v to types.Date", v)

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
func (p *DateList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []Date
	list = make([]Date, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = Date{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

// datetime
type DateTime struct {
	time.Time
}

func NowDateTime() DateTime {
	return DateTime{time.Now()}
}

func NowDateTimeP() *DateTime {
	return &DateTime{time.Now()}
}

func (t DateTime) ToString() string {
	return t.Format(`2006-01-02 15:04:05`)
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*t = DateTime{Time: now}
	return err
}

// Value insert timestamp into mysql need this function.
func (t DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof jstime.Time
func (t *DateTime) Scan(v any) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime{value}
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
	return time.Unix(t.Unix(), 0)
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
func (p *DateTimeList) Scan(data any) error {
	array := pgtype.TimestampArray{}
	err := array.Scan(data)
	if err != nil {
		return err
	}
	var list []DateTime
	list = make([]DateTime, len(array.Elements))
	for i, element := range array.Elements {
		list[i] = DateTime{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

func (d Date) AddTime(t Time) DateTime {
	return DateTime{time.Date(
		d.Time.Year(),
		d.Time.Month(),
		d.Time.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(), 0, nil,
	)}
}

func (t Time) AddData(d Date) DateTime {
	return DateTime{time.Date(
		d.Time.Year(),
		d.Time.Month(),
		d.Time.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(), 0, nil,
	)}
}

// datetime
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
func (t *AutoDateTime) Scan(v any) error {
	var s = ""
	switch v := v.(type) {
	case string:
		s = v[:8]
	case []byte:
		s = string(v)[:8]
	case time.Time:
		*t = AutoDateTime{v}
	case Time:
		*t = AutoDateTime{v.Time}
	case Date:
		*t = AutoDateTime{v.Time}
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
	*t = AutoDateTime{Time: now}
	return nil
}
