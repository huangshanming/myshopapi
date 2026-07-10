package server

import (
	"context"
	"fmt"
	"net"

	userv1 "mymall/api/gen/user/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/services/user-service/internal/service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	userv1.UnimplementedUserServiceServer
	svc *service.UserService
}

func NewUserGRPCServer(svc *service.UserService) *UserGRPCServer {
	return &UserGRPCServer{svc: svc}
}

func (s *UserGRPCServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := s.svc.GetProfile(req.GetUserId())
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

func Listen(port int, svc *service.UserService, logger *zap.Logger) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, err
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Logging(logger)))
	userv1.RegisterUserServiceServer(s, NewUserGRPCServer(svc))
	return s, lis, nil
}
