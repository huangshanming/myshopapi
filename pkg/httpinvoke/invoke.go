package httpinvoke

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/zeromicro/go-zero/rest/pathvar"

	"mymall/pkg/xerr"
)

// Run calls a legacy (w,r) handler with a synthetic request and returns the JSON body.
// On HTTP 4xx/5xx it returns xerr with the status and msg from the body when present.
func Run(ctx context.Context, method, path string, pathVars map[string]string, query url.Values, body any, h http.HandlerFunc) (json.RawMessage, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	if query == nil {
		query = url.Values{}
	}
	u := path
	if qs := query.Encode(); qs != "" {
		u = path + "?" + qs
	}
	req := httptest.NewRequest(method, u, rdr)
	req = req.WithContext(ctx)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if len(pathVars) > 0 {
		req = pathvar.WithVars(req, pathVars)
	}
	rr := httptest.NewRecorder()
	h(rr, req)
	if rr.Code >= 400 {
		var er struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &er)
		msg := er.Msg
		if msg == "" {
			msg = string(bytes.TrimSpace(rr.Body.Bytes()))
		}
		if msg == "" {
			msg = http.StatusText(rr.Code)
		}
		code := er.Code
		if code == 0 {
			code = rr.Code
		}
		return nil, xerr.New(code, msg)
	}
	if rr.Body.Len() == 0 {
		return nil, nil
	}
	out := make(json.RawMessage, rr.Body.Len())
	copy(out, rr.Body.Bytes())
	return out, nil
}

// Decode unmarshals raw into dest; nil raw is OK.
func Decode(raw json.RawMessage, dest any) error {
	if len(raw) == 0 || dest == nil {
		return nil
	}
	return json.Unmarshal(raw, dest)
}
