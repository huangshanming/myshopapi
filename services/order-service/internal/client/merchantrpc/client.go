package merchantrpc

import (
	"context"
	"fmt"

	merchantv1 "mymall/api/gen/merchant/v1"
	"mymall/api/rpcclient/merchantservice"
	"mymall/pkg/zrpcx"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/status"
)

type Client struct {
	cli    zrpc.Client
	client merchantservice.MerchantService
}

type ConsumeResult struct {
	EntryID      uint64
	ProductID    uint64
	SeckillPrice float64
	Quantity     int
}

type MatchItem struct {
	ProductID      uint64
	CategoryID     uint64
	Amount         float64
	SeckillEntryID uint64
}

type MatchCouponView struct {
	UserCouponID    uint64
	CouponID        uint64
	Name            string
	CouponType      string
	DiscountAmount  float64
	ThresholdAmount float64
	ValidEnd        string
	Usable          bool
	Reason          string
	Best            bool
}

type MatchResp struct {
	GoodsAmount      float64
	DiscountAmount   float64
	PayAmount        float64
	BestUserCouponID uint64
	Available        []MatchCouponView
	Unavailable      []MatchCouponView
}

func New(addr string, etcdHosts []string) (*Client, error) {
	cli, err := zrpcx.Dial(addr, etcdHosts, zrpcx.KeyMerchant)
	if err != nil {
		return nil, fmt.Errorf("merchant zrpc dial: %w", err)
	}
	return &Client{
		cli:    cli,
		client: merchantservice.NewMerchantService(cli),
	}, nil
}

func (c *Client) Close() error {
	if c != nil && c.cli != nil {
		c.cli.Conn().Close()
	}
	return nil
}

func (c *Client) Consume(ctx context.Context, entryID, productID uint64, qty int) (*ConsumeResult, error) {
	if c == nil {
		return nil, fmt.Errorf("商户服务不可用")
	}
	resp, err := c.client.ConsumeSeckill(ctx, &merchantv1.ConsumeSeckillRequest{
		EntryId: entryID, ProductId: productID, Quantity: int32(qty),
	})
	if err != nil {
		return nil, rpcErr(err)
	}
	return &ConsumeResult{
		EntryID: resp.GetEntryId(), ProductID: resp.GetProductId(),
		SeckillPrice: resp.GetSeckillPrice(), Quantity: int(resp.GetQuantity()),
	}, nil
}

func (c *Client) Restore(ctx context.Context, entryID uint64, qty int) error {
	if c == nil || entryID == 0 || qty < 1 {
		return nil
	}
	_, err := c.client.RestoreSeckill(ctx, &merchantv1.RestoreSeckillRequest{
		EntryId: entryID, Quantity: int32(qty),
	})
	return rpcErr(err)
}

func (c *Client) MatchCoupons(ctx context.Context, userID, shopID, userCouponID uint64, items []MatchItem) (*MatchResp, error) {
	if c == nil {
		return nil, fmt.Errorf("商户服务不可用")
	}
	pbItems := make([]*merchantv1.MatchItem, 0, len(items))
	for _, it := range items {
		pbItems = append(pbItems, &merchantv1.MatchItem{
			ProductId: it.ProductID, CategoryId: it.CategoryID,
			Amount: it.Amount, SeckillEntryId: it.SeckillEntryID,
		})
	}
	resp, err := c.client.MatchCoupons(ctx, &merchantv1.MatchCouponsRequest{
		UserId: userID, ShopId: shopID, UserCouponId: userCouponID, Items: pbItems,
	})
	if err != nil {
		return nil, rpcErr(err)
	}
	out := &MatchResp{
		GoodsAmount: resp.GetGoodsAmount(), DiscountAmount: resp.GetDiscountAmount(),
		PayAmount: resp.GetPayAmount(), BestUserCouponID: resp.GetBestUserCouponId(),
		Available: []MatchCouponView{}, Unavailable: []MatchCouponView{},
	}
	for _, v := range resp.GetAvailable() {
		out.Available = append(out.Available, fromView(v))
	}
	for _, v := range resp.GetUnavailable() {
		out.Unavailable = append(out.Unavailable, fromView(v))
	}
	return out, nil
}

func fromView(v *merchantv1.MatchCouponView) MatchCouponView {
	if v == nil {
		return MatchCouponView{}
	}
	return MatchCouponView{
		UserCouponID: v.GetUserCouponId(), CouponID: v.GetCouponId(), Name: v.GetName(),
		CouponType: v.GetCouponType(), DiscountAmount: v.GetDiscountAmount(),
		ThresholdAmount: v.GetThresholdAmount(), ValidEnd: v.GetValidEnd(),
		Usable: v.GetUsable(), Reason: v.GetReason(), Best: v.GetBest(),
	}
}

func (c *Client) LockCoupon(ctx context.Context, userCouponID, userID, orderID uint64, discount float64) error {
	if c == nil {
		return fmt.Errorf("商户服务不可用")
	}
	_, err := c.client.LockCoupon(ctx, &merchantv1.LockCouponRequest{
		UserCouponId: userCouponID, UserId: userID, OrderId: orderID, DiscountAmount: discount,
	})
	return rpcErr(err)
}

func (c *Client) UnlockCoupon(ctx context.Context, userCouponID, orderID uint64) error {
	if c == nil || userCouponID == 0 {
		return nil
	}
	_, err := c.client.UnlockCoupon(ctx, &merchantv1.UnlockCouponRequest{
		UserCouponId: userCouponID, OrderId: orderID,
	})
	return rpcErr(err)
}

func (c *Client) RedeemCoupon(ctx context.Context, userCouponID, orderID uint64, discount float64) error {
	if c == nil || userCouponID == 0 {
		return nil
	}
	_, err := c.client.RedeemCoupon(ctx, &merchantv1.RedeemCouponRequest{
		UserCouponId: userCouponID, OrderId: orderID, DiscountAmount: discount,
	})
	return rpcErr(err)
}

func (c *Client) ReturnCoupon(ctx context.Context, userCouponID, orderID uint64) error {
	if c == nil || userCouponID == 0 {
		return nil
	}
	_, err := c.client.ReturnCoupon(ctx, &merchantv1.ReturnCouponRequest{
		UserCouponId: userCouponID, OrderId: orderID,
	})
	return rpcErr(err)
}

func (c *Client) OrderGiftCoupons(ctx context.Context, userID, shopID uint64) error {
	if c == nil {
		return fmt.Errorf("商户服务不可用")
	}
	_, err := c.client.OrderGiftCoupons(ctx, &merchantv1.OrderGiftCouponsRequest{
		UserId: userID, ShopId: shopID,
	})
	return rpcErr(err)
}

func rpcErr(err error) error {
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok && st.Message() != "" {
		return fmt.Errorf("%s", st.Message())
	}
	return fmt.Errorf("商户服务错误: %w", err)
}
