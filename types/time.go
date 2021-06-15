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

func (t Time) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"15:04:05"`)
	return []byte(tune), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"15:04:05"`, string(data), time.Local)
	*t = Time{Time: now}
	return err
}

func (t Time) ToGoTime() time.Time {
	return time.Unix(t.Unix(), 0)
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof jstime.Time
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
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
	s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
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
		list[i] = Time{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

//date
type Date struct {
	time.Time
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
func (t *Date) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Date{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
func (t Date) ToGoTime() time.Time {
	return time.Unix(t.Unix(), 0)
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
	s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
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
		list[i] = Date{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

//datetime
type DateTime struct {
	time.Time
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
func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
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
	s = s[:0] + "{" + s[1:len(s)-1] + "}" + s[len(s):]
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
		list[i] = DateTime{element.Time}
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &p)
	return err
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (t NullTime) IsNull() bool {
	return !t.Valid
}

// Scan implements the Scanner interface.
func (t *NullTime) Scan(v interface{}) error {
	if v == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	t.Valid = true
	value, ok := v.(time.Time)
	if ok {
		t.Time = value
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Value implements the driver Valuer interface.
func (t NullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}
	tune := t.Time.Format(`"15:04:05"`)
	return []byte(tune), nil
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = NullTime{
			Time:  time.Time{},
			Valid: false,
		}
		return nil
	}
	now, err := time.ParseInLocation(`"15:04:05"`, string(data), time.Local)
	if err != nil {
		*t = NullTime{
			Time:  now,
			Valid: true,
		}
		return nil
	} else {
		*t = NullTime{
			Time:  time.Time{},
			Valid: false,
		}
		return nil
	}
}

func (t *NullTime) SetTime(v time.Time) {
	t.Time = v
	t.Valid = true
}

type NullDateTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (t NullDateTime) IsNull() bool {
	return !t.Valid
}

// Scan implements the Scanner interface.
func (t *NullDateTime) Scan(v interface{}) error {
	if v == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	t.Valid = true
	value, ok := v.(time.Time)
	if ok {
		t.Time = value
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Value implements the driver Valuer interface.
func (t NullDateTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t NullDateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}
	tune := t.Time.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

func (t *NullDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = NullDateTime{
			Time:  time.Time{},
			Valid: false,
		}
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	if err != nil {
		*t = NullDateTime{
			Time:  now,
			Valid: true,
		}
		return nil
	} else {
		*t = NullDateTime{
			Time:  time.Time{},
			Valid: false,
		}
		return nil
	}
}

func (t *NullDateTime) SetTime(v time.Time) {
	t.Time = v
	t.Valid = true
}

type NullDate struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (t NullDate) IsNull() bool {
	return !t.Valid
}

// Scan implements the Scanner interface.
func (t *NullDate) Scan(v interface{}) error {
	if v == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}
	t.Valid = true
	value, ok := v.(time.Time)
	if ok {
		t.Time = value
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Value implements the driver Valuer interface.
func (t NullDate) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t NullDate) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil)
	}
	tune := t.Time.Format(`"2006-01-02"`)
	return []byte(tune), nil
}

func (t *NullDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = NullDate{
			Time:  time.Time{},
			Valid: false,
		}
		return nil
	}

	now, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	if err != nil {
		*t = NullDate{
			Time:  now,
			Valid: true,
		}
		return nil
	} else {
		*t = NullDate{
			Time:  time.Time{},
			Valid: false,
		}
		return nil
	}
}

func (t *NullDate) SetTime(v time.Time) {
	t.Time = v
	t.Valid = true
}
