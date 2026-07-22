package server

import (
	"fmt"

	catalogv1 "mymall/api/gen/catalog/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/pkg/zrpcx"
	"mymall/services/catalog-service/internal/svc"

	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StartZRPC(port int, etcdHosts []string, svcCtx *svc.ServiceContext, logger *zap.Logger) *zrpc.RpcServer {
	c := zrpcx.ServerConf(fmt.Sprintf("0.0.0.0:%d", port), etcdHosts, zrpcx.KeyCatalog)
	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		catalogv1.RegisterCatalogServiceServer(grpcServer, NewCatalogServiceServer(svcCtx))
	})
	s.AddUnaryInterceptors(interceptor.Logging(logger))
	return s
}
