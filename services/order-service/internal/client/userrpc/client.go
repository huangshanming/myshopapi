package userrpc

import (
	"context"
	"fmt"

	userv1 "mymall/api/gen/user/v1"

	"github.com/zeromicro/go-zero/zrpc"
)

type Client struct {
	cli    zrpc.Client
	client userv1.UserServiceClient
}

func New(addr string) (*Client, error) {
	c := zrpc.RpcClientConf{
		Endpoints: []string{addr},
		NonBlock:  true,
	}
	cli, err := zrpc.NewClient(c)
	if err != nil {
		return nil, fmt.Errorf("user zrpc dial: %w", err)
	}
	return &Client{
		cli:    cli,
		client: userv1.NewUserServiceClient(cli.Conn()),
	}, nil
}

func (c *Client) Close() error {
	if c.cli != nil {
		c.cli.Conn().Close()
	}
	return nil
}

func (c *Client) GetUser(ctx context.Context, userID uint64) (*userv1.GetUserResponse, error) {
	return c.client.GetUser(ctx, &userv1.GetUserRequest{UserId: userID})
}
