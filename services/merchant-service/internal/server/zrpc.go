package server

import (
	"fmt"

	merchantv1 "mymall/api/gen/merchant/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/pkg/zrpcx"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StartZRPC(port int, etcdHosts []string, svcCtx *svc.ServiceContext, logger *zap.Logger) *zrpc.RpcServer {
	c := zrpcx.ServerConf(fmt.Sprintf("0.0.0.0:%d", port), etcdHosts, zrpcx.KeyMerchant)
	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		merchantv1.RegisterMerchantServiceServer(grpcServer, NewMerchantServiceServer(svcCtx))
	})
	s.AddUnaryInterceptors(interceptor.Logging(logger))
	return s
}
