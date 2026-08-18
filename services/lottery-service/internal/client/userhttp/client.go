package userhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string, timeoutSec int) *Client {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = "http://127.0.0.1:8881"
	}
	if timeoutSec <= 0 {
		timeoutSec = 8
	}
	return &Client{
		baseURL: base,
		http:    &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

type LedgerReq struct {
	UserID     uint64 `json:"user_id"`
	Points     int64  `json:"points"`
	Reason     string `json:"reason,optional"`
	RefNo      string `json:"ref_no,optional"`
	ChangeType string `json:"change_type,optional"`
	RefType    string `json:"ref_type,optional"`
}

type Address struct {
	ID            uint64 `json:"id"`
	UserID        uint64 `json:"user_id"`
	ReceiverName  string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	Detail        string `json:"detail"`
}

func (a *Address) FullAddress() string {
	if a == nil {
		return ""
	}
	return a.Province + a.City + a.District + a.Detail
}

type NotifyReq struct {
	UserID  uint64 `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type,omitempty"`
}

type apiErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c *Client) Deduct(ctx context.Context, req LedgerReq) error {
	return c.post(ctx, "/api/v1/internal/points/deduct", req)
}

func (c *Client) Refund(ctx context.Context, req LedgerReq) error {
	return c.post(ctx, "/api/v1/internal/points/refund", req)
}

func (c *Client) Add(ctx context.Context, req LedgerReq) error {
	return c.post(ctx, "/api/v1/internal/points/add", req)
}

func (c *Client) Notify(ctx context.Context, req NotifyReq) error {
	if req.Type == "" {
		req.Type = "lottery"
	}
	return c.post(ctx, "/api/v1/internal/notifications", req)
}

func (c *Client) GetAddress(ctx context.Context, userID, addressID uint64) (*Address, error) {
	q := url.Values{}
	q.Set("id", strconv.FormatUint(addressID, 10))
	q.Set("user_id", strconv.FormatUint(userID, 10))
	raw, err := c.get(ctx, "/api/v1/user/addresses/internal?"+q.Encode())
	if err != nil {
		return nil, err
	}
	var wrap struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, fmt.Errorf("解析地址响应失败")
	}
	payload := wrap.Data
	if len(payload) == 0 {
		payload = raw
	}
	var addr Address
	if err := json.Unmarshal(payload, &addr); err != nil {
		return nil, fmt.Errorf("解析地址失败")
	}
	if addr.ID == 0 {
		return nil, fmt.Errorf("收货地址无效")
	}
	return &addr, nil
}

func (c *Client) post(ctx context.Context, path string, body any) error {
	raw, _ := json.Marshal(body)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return fmt.Errorf("调用用户服务失败: %w", err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("%s", errMsg(data, resp.StatusCode))
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用用户服务失败: %w", err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return data, nil
	}
	return nil, fmt.Errorf("%s", errMsg(data, resp.StatusCode))
}

func errMsg(data []byte, status int) string {
	var ae apiErr
	_ = json.Unmarshal(data, &ae)
	msg := strings.TrimSpace(ae.Msg)
	if msg == "" {
		msg = strings.TrimSpace(string(data))
	}
	if msg == "" {
		msg = fmt.Sprintf("用户服务错误 HTTP %d", status)
	}
	return msg
}
