package merchanthttp

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

type ConsumeResult struct {
	EntryID      uint64  `json:"entry_id"`
	ProductID    uint64  `json:"product_id"`
	SeckillPrice float64 `json:"seckill_price"`
	Quantity     int     `json:"quantity"`
}

func New(base string) *Client {
	if base == "" {
		base = os.Getenv("MYMALL_MERCHANT_HTTP")
	}
	if base == "" {
		base = "http://127.0.0.1:8884"
	}
	return &Client{
		base: base,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Consume(ctx context.Context, entryID, productID uint64, qty int) (*ConsumeResult, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"entry_id":   entryID,
		"product_id": productID,
		"quantity":   qty,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/seckill/consume", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("秒杀服务不可用: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	var wrap apiResp
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, errorsNew("秒杀服务响应异常")
	}
	if wrap.Code != 200 {
		if wrap.Msg != "" {
			return nil, errorsNew(wrap.Msg)
		}
		return nil, errorsNew("秒杀扣库存失败")
	}
	var out ConsumeResult
	if err := json.Unmarshal(wrap.Data, &out); err != nil {
		return nil, errorsNew("秒杀服务响应异常")
	}
	return &out, nil
}

func (c *Client) Restore(ctx context.Context, entryID uint64, qty int) error {
	if entryID == 0 || qty < 1 {
		return nil
	}
	body, _ := json.Marshal(map[string]interface{}{
		"entry_id": entryID,
		"quantity": qty,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/seckill/restore", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	var wrap apiResp
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return err
	}
	if wrap.Code != 200 {
		return errorsNew(wrap.Msg)
	}
	return nil
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }

func errorsNew(s string) error { return simpleErr(s) }
