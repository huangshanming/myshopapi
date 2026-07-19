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
	var wrap apiResp
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, simpleErr("用户地址服务响应异常")
	}
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

type simpleErr string

func (e simpleErr) Error() string { return string(e) }
