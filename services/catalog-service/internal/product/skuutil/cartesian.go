package skuutil

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// SpecItem 规格定义
type SpecItem struct {
	Name   string
	Values []string
}

// Cartesian 根据规格项生成全部组合；每项为 name->value 映射
func Cartesian(specs []SpecItem) []map[string]string {
	if len(specs) == 0 {
		return []map[string]string{{}}
	}
	var result []map[string]string
	var walk func(int, map[string]string)
	walk = func(idx int, cur map[string]string) {
		if idx == len(specs) {
			cp := make(map[string]string, len(cur))
			for k, v := range cur {
				cp[k] = v
			}
			result = append(result, cp)
			return
		}
		sp := specs[idx]
		vals := sp.Values
		if len(vals) == 0 {
			walk(idx+1, cur)
			return
		}
		for _, v := range vals {
			cur[sp.Name] = v
			walk(idx+1, cur)
		}
	}
	walk(0, map[string]string{})
	return result
}

// SpecKey 规范化唯一键（排序后拼接，避免顺序导致重复）
func SpecKey(values map[string]string) string {
	if len(values) == 0 {
		return "default"
	}
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+values[k])
	}
	raw := strings.Join(parts, "|")
	sum := md5.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func SpecValuesJSON(values map[string]string) (string, error) {
	if values == nil {
		values = map[string]string{}
	}
	b, err := json.Marshal(values)
	return string(b), err
}

func SkuNo(productNo string, specKey string) string {
	short := specKey
	if len(short) > 8 {
		short = short[:8]
	}
	return fmt.Sprintf("%s-%s", productNo, short)
}
