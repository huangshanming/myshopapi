package server

import (
	"context"
	"fmt"

	userv1 "mymall/api/gen/user/v1"
	"mymall/pkg/grpc/interceptor"
	biz "mymall/services/user-service/internal/biz"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	logic *biz.UserLogic
}

func NewUserServer(l *biz.UserLogic) *UserServer {
	return &UserServer{logic: l}
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := s.logic.GetProfile(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &userv1.GetUserResponse{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Status:   int32(user.Status),
	}, nil
}

func StartZRPC(port int, l *biz.UserLogic, logger *zap.Logger) *zrpc.RpcServer {
	c := zrpc.RpcServerConf{
		ListenOn: fmt.Sprintf("0.0.0.0:%d", port),
	}
	c.Mode = service.DevMode
	c.Log.Mode = "console"
	c.Log.Encoding = "plain"

	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		userv1.RegisterUserServiceServer(grpcServer, NewUserServer(l))
	})
	s.AddUnaryInterceptors(interceptor.Logging(logger))
	return s
}
