package userhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	base   string
	client *http.Client
}

type apiResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func New(base string) *Client {
	if base == "" {
		base = os.Getenv("MYMALL_USER_HTTP")
	}
	if base == "" {
		base = "http://127.0.0.1:8881"
	}
	return &Client{
		base:   base,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

type NotifyReq struct {
	UserID   uint64 `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	MsgType  string `json:"msg_type"`
	LinkType string `json:"link_type"`
	LinkID   uint64 `json:"link_id"`
	Extra    string `json:"extra"`
}

func (c *Client) Notify(ctx context.Context, req NotifyReq) error {
	if c == nil || req.UserID == 0 || req.Title == "" {
		return nil
	}
	if req.MsgType == "" {
		req.MsgType = "system"
	}
	if req.LinkType == "" {
		req.LinkType = "none"
	}
	if req.Extra == "" {
		req.Extra = "{}"
	}
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/internal/notifications", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("通知服务不可用: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		var wrap apiResp
		if json.Unmarshal(raw, &wrap) == nil && wrap.Msg != "" {
			return fmt.Errorf("%s", wrap.Msg)
		}
		return fmt.Errorf("通知失败(%d)", res.StatusCode)
	}
	return nil
}
