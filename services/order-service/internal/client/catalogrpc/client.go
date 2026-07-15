package catalogrpc

import (
	"context"
	"fmt"

	catalogv1 "mymall/api/gen/catalog/v1"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Client struct {
	cli    zrpc.Client
	client catalogv1.CatalogServiceClient
}

func New(addr string) (*Client, error) {
	c := zrpc.RpcClientConf{
		Endpoints: []string{addr},
		NonBlock:  true, // 滚动发布时 catalog 可能尚未就绪，勿在启动时硬失败
	}
	cli, err := zrpc.NewClient(c)
	if err != nil {
		return nil, fmt.Errorf("catalog zrpc dial: %w", err)
	}
	return &Client{
		cli:    cli,
		client: catalogv1.NewCatalogServiceClient(cli.Conn()),
	}, nil
}

func (c *Client) Close() error {
	if c.cli != nil {
		c.cli.Conn().Close()
	}
	return nil
}

func (c *Client) Conn() grpc.ClientConnInterface {
	return c.cli.Conn()
}

func (c *Client) BatchGetProducts(ctx context.Context, ids []uint64) (*catalogv1.BatchGetProductsResponse, error) {
	return c.client.BatchGetProducts(ctx, &catalogv1.BatchGetProductsRequest{ProductIds: ids})
}
