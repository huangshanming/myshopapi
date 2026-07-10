package client

import (
	"context"
	"fmt"

	catalogv1 "mymall/api/gen/catalog/v1"
	userv1 "mymall/api/gen/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CatalogClient struct {
	conn   *grpc.ClientConn
	client catalogv1.CatalogServiceClient
}

func NewCatalogClient(addr string) (*CatalogClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("catalog grpc dial: %w", err)
	}
	return &CatalogClient{conn: conn, client: catalogv1.NewCatalogServiceClient(conn)}, nil
}

func (c *CatalogClient) Close() error {
	return c.conn.Close()
}

func (c *CatalogClient) BatchGetProducts(ctx context.Context, ids []uint64) (*catalogv1.BatchGetProductsResponse, error) {
	return c.client.BatchGetProducts(ctx, &catalogv1.BatchGetProductsRequest{ProductIds: ids})
}

type UserClient struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("user grpc dial: %w", err)
	}
	return &UserClient{conn: conn, client: userv1.NewUserServiceClient(conn)}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) GetUser(ctx context.Context, userID uint64) (*userv1.GetUserResponse, error) {
	return c.client.GetUser(ctx, &userv1.GetUserRequest{UserId: userID})
}
