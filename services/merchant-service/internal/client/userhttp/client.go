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
		base: base,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type PointsLedgerReq struct {
	UserID     uint64 `json:"user_id"`
	Points     int    `json:"points"`
	ChangeType string `json:"change_type"`
	Remark     string `json:"remark"`
	RefType    string `json:"ref_type"`
	RefID      uint64 `json:"ref_id"`
}

func (c *Client) DeductPoints(ctx context.Context, req PointsLedgerReq) error {
	return c.postLedger(ctx, "/api/v1/internal/points/deduct", req)
}

func (c *Client) RefundPoints(ctx context.Context, req PointsLedgerReq) error {
	return c.postLedger(ctx, "/api/v1/internal/points/refund", req)
}

func (c *Client) postLedger(ctx context.Context, path string, req PointsLedgerReq) error {
	if c == nil {
		return simpleErr("积分服务未配置")
	}
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("积分服务不可用: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	var wrap apiResp
	if err := json.Unmarshal(raw, &wrap); err != nil {
		if res.StatusCode >= 400 {
			return simpleErr("积分操作失败")
		}
		return nil
	}
	if wrap.Code != 0 && wrap.Code != 200 {
		if wrap.Msg != "" {
			return simpleErr(wrap.Msg)
		}
		return simpleErr("积分操作失败")
	}
	if res.StatusCode >= 400 {
		if wrap.Msg != "" {
			return simpleErr(wrap.Msg)
		}
		return simpleErr("积分操作失败")
	}
	return nil
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }
