package catalogrpc

import (
	"context"
	"fmt"

	catalogv1 "mymall/api/gen/catalog/v1"
	"mymall/api/rpcclient/catalogservice"
	"mymall/pkg/zrpcx"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Client struct {
	cli    zrpc.Client
	client catalogservice.CatalogService
}

func New(addr string, etcdHosts []string) (*Client, error) {
	cli, err := zrpcx.Dial(addr, etcdHosts, zrpcx.KeyCatalog)
	if err != nil {
		return nil, fmt.Errorf("catalog zrpc dial: %w", err)
	}
	return &Client{
		cli:    cli,
		client: catalogservice.NewCatalogService(cli),
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
