package appinput

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"mymall/pkg/pagination"
)

// CallInput carries parsed HTTP inputs without requiring handlers to use *http.Request.
// Request is optional and only used for multipart upload bridges.
type CallInput struct {
	Body     any
	PathVars map[string]string
	Query    url.Values
	Request  *http.Request
}

func (in CallInput) Path(key string) string {
	if in.PathVars == nil {
		return ""
	}
	return in.PathVars[key]
}

func (in CallInput) PathUint64(key string) uint64 {
	v, _ := strconv.ParseUint(in.Path(key), 10, 64)
	return v
}

func (in CallInput) QueryGet(key string) string {
	if in.Query == nil {
		return ""
	}
	return in.Query.Get(key)
}

func (in CallInput) Page() (page, pageSize int) {
	req := pagination.PageReq{Page: 1, PageSize: 10}
	if in.Query != nil {
		if p := in.Query.Get("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil {
				req.Page = n
			}
		}
		if p := in.Query.Get("page_size"); p != "" {
			if n, err := strconv.Atoi(p); err == nil {
				req.PageSize = n
			}
		}
	}
	page, pageSize, _ = pagination.Normalize(&req)
	return page, pageSize
}

func BindBody(in CallInput, dest any) error {
	if in.Body == nil {
		return nil
	}
	b, err := json.Marshal(in.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}
