package common

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LocalTime 自定义时间类型，用于格式化 JSON 输出
type LocalTime time.Time

const TimeFormat = "2006-01-02 15:04:05"

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", tt.Format(TimeFormat))
	return []byte(formatted), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = LocalTime(time.Time{})
		return nil
	}
	var tt time.Time
	if err := tt.UnmarshalJSON(data); err != nil {
		return err
	}
	*t = LocalTime(tt)
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if v == nil {
		*t = LocalTime(time.Time{})
		return nil
	}
	switch val := v.(type) {
	case time.Time:
		*t = LocalTime(val)
	case []byte:
		tt, err := time.Parse(TimeFormat, string(val))
		if err != nil {
			return err
		}
		*t = LocalTime(tt)
	case string:
		tt, err := time.Parse(TimeFormat, val)
		if err != nil {
			return err
		}
		*t = LocalTime(tt)
	default:
		return fmt.Errorf("不支持的类型：%T", v)
	}
	return nil
}
