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

type MatchItem struct {
	ProductID      uint64  `json:"product_id"`
	CategoryID     uint64  `json:"category_id"`
	Amount         float64 `json:"amount"`
	SeckillEntryID uint64  `json:"seckill_entry_id"`
}

type MatchCouponView struct {
	UserCouponID    uint64  `json:"user_coupon_id"`
	CouponID        uint64  `json:"coupon_id"`
	Name            string  `json:"name"`
	CouponType      string  `json:"coupon_type"`
	DiscountAmount  float64 `json:"discount_amount"`
	ThresholdAmount float64 `json:"threshold_amount"`
	ValidEnd        string  `json:"valid_end"`
	Usable          bool    `json:"usable"`
	Reason          string  `json:"reason,omitempty"`
	Best            bool    `json:"best,omitempty"`
}

type MatchResp struct {
	GoodsAmount      float64           `json:"goods_amount"`
	DiscountAmount   float64           `json:"discount_amount"`
	PayAmount        float64           `json:"pay_amount"`
	BestUserCouponID uint64            `json:"best_user_coupon_id"`
	Available        []MatchCouponView `json:"available"`
	Unavailable      []MatchCouponView `json:"unavailable"`
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

func (c *Client) post(ctx context.Context, path string, body interface{}, out interface{}) error {
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("商户服务不可用: %w", err)
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		var wrap apiResp
		if json.Unmarshal(respBody, &wrap) == nil && wrap.Msg != "" {
			return errorsNew(wrap.Msg)
		}
		return errorsNew(fmt.Sprintf("商户服务错误(%d)", res.StatusCode))
	}
	if out == nil || len(respBody) == 0 || string(respBody) == "null" {
		return nil
	}
	// 兼容旧信封
	var wrap apiResp
	if json.Unmarshal(respBody, &wrap) == nil && wrap.Code != 0 {
		if wrap.Code != 200 {
			if wrap.Msg != "" {
				return errorsNew(wrap.Msg)
			}
			return errorsNew("商户服务失败")
		}
		if wrap.Data != nil && out != nil {
			return json.Unmarshal(wrap.Data, out)
		}
		return nil
	}
	if out != nil {
		return json.Unmarshal(respBody, out)
	}
	return nil
}

func (c *Client) Consume(ctx context.Context, entryID, productID uint64, qty int) (*ConsumeResult, error) {
	var out ConsumeResult
	err := c.post(ctx, "/api/v1/seckill/consume", map[string]interface{}{
		"entry_id": entryID, "product_id": productID, "quantity": qty,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Restore(ctx context.Context, entryID uint64, qty int) error {
	if entryID == 0 || qty < 1 {
		return nil
	}
	return c.post(ctx, "/api/v1/seckill/restore", map[string]interface{}{
		"entry_id": entryID, "quantity": qty,
	}, nil)
}

func (c *Client) MatchCoupons(ctx context.Context, userID, shopID, userCouponID uint64, items []MatchItem) (*MatchResp, error) {
	var out MatchResp
	err := c.post(ctx, "/api/v1/internal/coupons/match", map[string]interface{}{
		"user_id": userID, "shop_id": shopID, "user_coupon_id": userCouponID, "items": items,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) LockCoupon(ctx context.Context, userCouponID, userID, orderID uint64, discount float64) error {
	return c.post(ctx, "/api/v1/internal/coupons/lock", map[string]interface{}{
		"user_coupon_id": userCouponID, "user_id": userID, "order_id": orderID, "discount_amount": discount,
	}, nil)
}

func (c *Client) UnlockCoupon(ctx context.Context, userCouponID, orderID uint64) error {
	if userCouponID == 0 {
		return nil
	}
	return c.post(ctx, "/api/v1/internal/coupons/unlock", map[string]interface{}{
		"user_coupon_id": userCouponID, "order_id": orderID,
	}, nil)
}

func (c *Client) RedeemCoupon(ctx context.Context, userCouponID, orderID uint64, discount float64) error {
	if userCouponID == 0 {
		return nil
	}
	return c.post(ctx, "/api/v1/internal/coupons/redeem", map[string]interface{}{
		"user_coupon_id": userCouponID, "order_id": orderID, "discount_amount": discount,
	}, nil)
}

func (c *Client) ReturnCoupon(ctx context.Context, userCouponID, orderID uint64) error {
	if userCouponID == 0 {
		return nil
	}
	return c.post(ctx, "/api/v1/internal/coupons/return", map[string]interface{}{
		"user_coupon_id": userCouponID, "order_id": orderID,
	}, nil)
}

func (c *Client) OrderGiftCoupons(ctx context.Context, userID, shopID uint64) error {
	return c.post(ctx, "/api/v1/internal/coupons/order-gift", map[string]interface{}{
		"user_id": userID, "shop_id": shopID,
	}, nil)
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }

func errorsNew(s string) error { return simpleErr(s) }
