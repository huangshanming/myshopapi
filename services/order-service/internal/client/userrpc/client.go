package userrpc

import (
	"context"
	"fmt"

	userv1 "mymall/api/gen/user/v1"
	"mymall/api/rpcclient/userservice"
	"mymall/pkg/zrpcx"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/status"
)

type Client struct {
	cli    zrpc.Client
	client userservice.UserService
}

func New(addr string, etcdHosts []string) (*Client, error) {
	cli, err := zrpcx.Dial(addr, etcdHosts, zrpcx.KeyUser)
	if err != nil {
		return nil, fmt.Errorf("user zrpc dial: %w", err)
	}
	return &Client{
		cli:    cli,
		client: userservice.NewUserService(cli),
	}, nil
}

func (c *Client) Close() error {
	if c != nil && c.cli != nil {
		c.cli.Conn().Close()
	}
	return nil
}

func (c *Client) GetUser(ctx context.Context, userID uint64) (*userv1.GetUserResponse, error) {
	return c.client.GetUser(ctx, &userv1.GetUserRequest{UserId: userID})
}

type Address struct {
	ID            uint64
	UserID        uint64
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	Detail        string
	IsDefault     int
}

func (a *Address) FullAddress() string {
	if a == nil {
		return ""
	}
	return a.Province + a.City + a.District + a.Detail
}

func (c *Client) GetAddress(ctx context.Context, userID, addressID uint64) (*Address, error) {
	if c == nil {
		return nil, fmt.Errorf("用户地址服务不可用")
	}
	resp, err := c.client.GetAddress(ctx, &userv1.GetAddressRequest{UserId: userID, AddressId: addressID})
	if err != nil {
		return nil, rpcErr(err, "收货地址无效")
	}
	return &Address{
		ID:            resp.GetId(),
		UserID:        resp.GetUserId(),
		ReceiverName:  resp.GetReceiverName(),
		ReceiverPhone: resp.GetReceiverPhone(),
		Province:      resp.GetProvince(),
		City:          resp.GetCity(),
		District:      resp.GetDistrict(),
		Detail:        resp.GetDetail(),
		IsDefault:     int(resp.GetIsDefault()),
	}, nil
}

func (c *Client) Freeze(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	if c == nil {
		return fmt.Errorf("用户钱包服务不可用")
	}
	_, err := c.client.FreezeWallet(ctx, &userv1.WalletOpRequest{
		UserId: userID, Amount: amount, OrderId: orderID, OrderNo: orderNo,
	})
	return rpcErr(err, "钱包操作失败")
}

func (c *Client) Unfreeze(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	if c == nil {
		return fmt.Errorf("用户钱包服务不可用")
	}
	_, err := c.client.UnfreezeWallet(ctx, &userv1.WalletOpRequest{
		UserId: userID, Amount: amount, OrderId: orderID, OrderNo: orderNo,
	})
	return rpcErr(err, "钱包操作失败")
}

func (c *Client) Settle(ctx context.Context, userID uint64, amount float64, orderID uint64, orderNo string) error {
	if c == nil {
		return fmt.Errorf("用户钱包服务不可用")
	}
	_, err := c.client.SettleWallet(ctx, &userv1.WalletOpRequest{
		UserId: userID, Amount: amount, OrderId: orderID, OrderNo: orderNo,
	})
	return rpcErr(err, "钱包操作失败")
}

type NotifyReq struct {
	UserID   uint64
	Title    string
	Content  string
	MsgType  string
	LinkType string
	LinkID   uint64
	Extra    string
}

func (c *Client) Notify(ctx context.Context, req NotifyReq) error {
	if c == nil || req.UserID == 0 || req.Title == "" {
		return nil
	}
	if req.MsgType == "" {
		req.MsgType = "order"
	}
	if req.LinkType == "" {
		req.LinkType = "none"
	}
	_, err := c.client.Notify(ctx, &userv1.NotifyRequest{
		UserId: req.UserID, Title: req.Title, Content: req.Content,
		MsgType: req.MsgType, LinkType: req.LinkType, LinkId: req.LinkID, Extra: req.Extra,
	})
	return rpcErr(err, "通知失败")
}

type TaskEventReq struct {
	UserID   uint64
	TaskCode string
	Delta    int
	RefType  string
	RefID    uint64
}

func (c *Client) TaskEvent(ctx context.Context, req TaskEventReq) error {
	if c == nil || req.UserID == 0 || req.TaskCode == "" {
		return nil
	}
	delta := req.Delta
	if delta < 1 {
		delta = 1
	}
	_, err := c.client.TaskEvent(ctx, &userv1.TaskEventRequest{
		UserId: req.UserID, TaskCode: req.TaskCode, Delta: int32(delta),
		RefType: req.RefType, RefId: req.RefID,
	})
	_ = err // best-effort like HTTP client
	return nil
}

func (c *Client) GetConfig(ctx context.Context, key string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("用户配置服务不可用")
	}
	resp, err := c.client.GetConfig(ctx, &userv1.GetConfigRequest{Key: key})
	if err != nil {
		return "", rpcErr(err, "读取配置失败")
	}
	return resp.GetValue(), nil
}

func rpcErr(err error, fallback string) error {
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok && st.Message() != "" {
		return fmt.Errorf("%s", st.Message())
	}
	return fmt.Errorf("%s: %w", fallback, err)
}
