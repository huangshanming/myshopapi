package server

import (
	"fmt"

	userv1 "mymall/api/gen/user/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/pkg/zrpcx"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// StartZRPC starts the goctl-generated UserServiceServer with zrpcx (etcd optional).
func StartZRPC(port int, etcdHosts []string, svcCtx *svc.ServiceContext, logger *zap.Logger) *zrpc.RpcServer {
	c := zrpcx.ServerConf(fmt.Sprintf("0.0.0.0:%d", port), etcdHosts, zrpcx.KeyUser)
	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		userv1.RegisterUserServiceServer(grpcServer, NewUserServiceServer(svcCtx))
	})
	s.AddUnaryInterceptors(interceptor.Logging(logger))
	return s
}
