package types

import "reflect"

// FromPaged unwraps legacy biz maps / pagination.PageRes into the go-zero PageListResp
// contract ({total, list: [...]}) expected by mall-uni / admin-web.
func FromPaged(data any) *PageListResp {
	if data == nil {
		return &PageListResp{List: []any{}}
	}

	switch v := data.(type) {
	case map[string]interface{}:
		return &PageListResp{Total: asInt64(v["total"]), List: coalesceList(v["list"])}
	}

	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return &PageListResp{List: []any{}}
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return &PageListResp{Total: int64(rv.Len()), List: data}
	case reflect.Struct:
		out := &PageListResp{}
		if f := rv.FieldByName("List"); f.IsValid() && f.CanInterface() {
			out.List = coalesceList(f.Interface())
		}
		if f := rv.FieldByName("Total"); f.IsValid() && f.CanInterface() {
			out.Total = asInt64(f.Interface())
		}
		if out.List == nil {
			out.List = []any{}
		}
		return out
	default:
		return &PageListResp{List: data}
	}
}

func coalesceList(v any) any {
	if v == nil {
		return []any{}
	}
	return v
}

func asInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case uint64:
		return int64(n)
	case float64:
		return int64(n)
	case float32:
		return int64(n)
	default:
		return 0
	}
}
