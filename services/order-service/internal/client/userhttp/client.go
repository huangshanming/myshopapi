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

func (c *Client) postOp(ctx context.Context, path string, userID uint64, amount float64, orderID uint64, orderNo string) error {
	body, _ := json.Marshal(map[string]interface{}{
		"user_id":  userID,
		"amount":   amount,
		"order_id": orderID,
		"order_no": orderNo,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("用户钱包服务不可用: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	var wrap apiResp
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return simpleErr("用户钱包服务响应异常")
	}
	if wrap.Code != 200 {
		if wrap.Msg != "" {
			return simpleErr(wrap.Msg)
		}
		return simpleErr("钱包操作失败")
	}
	return nil
}

func (c *Client) Freeze(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	return c.postOp(ctx, "/api/v1/user/wallet/freeze", userID, amount, orderID, orderNo)
}

func (c *Client) Unfreeze(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	return c.postOp(ctx, "/api/v1/user/wallet/unfreeze", userID, amount, orderID, orderNo)
}

func (c *Client) Settle(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	return c.postOp(ctx, "/api/v1/user/wallet/settle", userID, amount, orderID, orderNo)
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
	IsDefault     int    `json:"is_default"`
}

func (a *Address) FullAddress() string {
	if a == nil {
		return ""
	}
	return a.Province + a.City + a.District + a.Detail
}

func (c *Client) GetAddress(ctx context.Context, userID, addressID uint64) (*Address, error) {
	url := fmt.Sprintf("%s/api/v1/user/addresses/internal?user_id=%d&id=%d", c.base, userID, addressID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("用户地址服务不可用: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	// 新契约：直接 DTO；旧信封兼容
	var wrap apiResp
	if json.Unmarshal(raw, &wrap) == nil && wrap.Code != 0 {
		if wrap.Code != 200 {
			if wrap.Msg != "" {
				return nil, simpleErr(wrap.Msg)
			}
			return nil, simpleErr("收货地址无效")
		}
		var out Address
		if err := json.Unmarshal(wrap.Data, &out); err != nil {
			return nil, simpleErr("用户地址服务响应异常")
		}
		return &out, nil
	}
	if res.StatusCode >= 400 {
		return nil, simpleErr("收货地址无效")
	}
	var out Address
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, simpleErr("用户地址服务响应异常")
	}
	return &out, nil
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
	if req.UserID == 0 || req.Title == "" {
		return nil
	}
	if req.MsgType == "" {
		req.MsgType = "order"
	}
	if req.LinkType == "" {
		req.LinkType = "none"
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
			return simpleErr(wrap.Msg)
		}
		return simpleErr(fmt.Sprintf("通知失败(%d)", res.StatusCode))
	}
	return nil
}

type TaskEventReq struct {
	UserID   uint64 `json:"user_id"`
	TaskCode string `json:"task_code"`
	Delta    int    `json:"delta"`
	RefType  string `json:"ref_type"`
	RefID    uint64 `json:"ref_id"`
}

func (c *Client) TaskEvent(ctx context.Context, req TaskEventReq) error {
	if c == nil || req.UserID == 0 || req.TaskCode == "" {
		return nil
	}
	if req.Delta < 1 {
		req.Delta = 1
	}
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/internal/tasks/events", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("任务服务不可用: %w", err)
	}
	defer res.Body.Close()
	io.Copy(io.Discard, res.Body)
	return nil
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }
