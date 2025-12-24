package common

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type PageReq struct {
	Page     int    `form:"page" json:"page" comment:"当前页码，默认1"`
	PageSize int    `form:"page_size" json:"page_size" comment:"每页条数，默认10，最大100"`
	OrderBy  string `form:"order_by" json:"order_by" comment:"排序字段，如id DESC"` // 可选：扩展排序参数
}

// PageRes 全局统一分页响应结果（支持任意数据类型的列表）
type PageRes[T any] struct {
	Total     int64 `json:"total" comment:"总记录数"`
	Page      int   `json:"page" comment:"当前页码"`
	PageSize  int   `json:"page_size" comment:"每页条数"`
	TotalPage int   `json:"total_page" comment:"总页数"`
	List      []T   `json:"list" comment:"当前页数据列表"`
}

// 自定义时间类型，用于格式化 JSON 输出
type LocalTime time.Time

// 目标时间格式模板
const TimeFormat = "2006-01-02 15:04:05"

// MarshalJSON 实现 json.Marshaler 接口，序列化时自动格式化
func (t LocalTime) MarshalJSON() ([]byte, error) {
	// 将自定义类型转为 time.Time
	tt := time.Time(t)
	// 若时间为零值，返回 null
	if tt.IsZero() {
		return []byte("null"), nil
	}
	// 格式化时间并包裹为 JSON 字符串
	formatted := fmt.Sprintf("\"%s\"", tt.Format(TimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，反序列化时解析时间
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	// 处理 null 值
	if string(data) == "null" {
		*t = LocalTime(time.Time{})
		return nil
	}
	// 解析 JSON 字符串为 time.Time
	var tt time.Time
	err := tt.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	*t = LocalTime(tt)
	return nil
}

// Value 实现 driver.Valuer 接口，适配 Gorm 写入数据库
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner 接口，适配 Gorm 从数据库读取
func (t *LocalTime) Scan(v interface{}) error {
	if v == nil {
		*t = LocalTime(time.Time{})
		return nil
	}
	switch val := v.(type) {
	case time.Time:
		*t = LocalTime(val)
	case []byte:
		tt, err := time.Parse("2006-01-02 15:04:05", string(val))
		if err != nil {
			return err
		}
		*t = LocalTime(tt)
	case string:
		tt, err := time.Parse("2006-01-02 15:04:05", val)
		if err != nil {
			return err
		}
		*t = LocalTime(tt)
	default:
		return fmt.Errorf("不支持的类型：%T", v)
	}
	return nil
}
