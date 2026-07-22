package userrpc

import (
	"context"
	"fmt"

	userv1 "mymall/api/gen/user/v1"
	"mymall/api/rpcclient/userservice"
	"mymall/pkg/zrpcx"

	"github.com/zeromicro/go-zero/zrpc"
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
		req.MsgType = "system"
	}
	if req.LinkType == "" {
		req.LinkType = "none"
	}
	if req.Extra == "" {
		req.Extra = "{}"
	}
	_, err := c.client.Notify(ctx, &userv1.NotifyRequest{
		UserId: req.UserID, Title: req.Title, Content: req.Content,
		MsgType: req.MsgType, LinkType: req.LinkType, LinkId: req.LinkID, Extra: req.Extra,
	})
	return err
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
	_, _ = c.client.TaskEvent(ctx, &userv1.TaskEventRequest{
		UserId: req.UserID, TaskCode: req.TaskCode, Delta: int32(delta),
		RefType: req.RefType, RefId: req.RefID,
	})
	return nil
}
