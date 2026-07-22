package server

import (
	"context"
	"fmt"

	userv1 "mymall/api/gen/user/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	svcCtx *svc.ServiceContext
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{svcCtx: svcCtx}
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := biz.NewUserLogic(s.svcCtx).GetProfile(ctx, req.GetUserId())
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

func (s *UserServer) GetAddress(ctx context.Context, req *userv1.GetAddressRequest) (*userv1.GetAddressResponse, error) {
	a, err := biz.NewAddressLogic(s.svcCtx).Get(ctx, req.GetUserId(), req.GetAddressId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}
	return &userv1.GetAddressResponse{
		Id:            a.ID,
		UserId:        a.UserID,
		ReceiverName:  a.ReceiverName,
		ReceiverPhone: a.ReceiverPhone,
		Province:      a.Province,
		City:          a.City,
		District:      a.District,
		Detail:        a.Detail,
		IsDefault:     int32(a.IsDefault),
	}, nil
}

func (s *UserServer) FreezeWallet(ctx context.Context, req *userv1.WalletOpRequest) (*userv1.EmptyResponse, error) {
	if err := biz.NewWalletLogic(s.svcCtx).FreezeForOrder(ctx, req.GetUserId(), req.GetAmount(), req.GetOrderId(), req.GetOrderNo()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}

func (s *UserServer) UnfreezeWallet(ctx context.Context, req *userv1.WalletOpRequest) (*userv1.EmptyResponse, error) {
	if err := biz.NewWalletLogic(s.svcCtx).UnfreezeOrder(ctx, req.GetUserId(), req.GetAmount(), req.GetOrderId(), req.GetOrderNo()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}

func (s *UserServer) SettleWallet(ctx context.Context, req *userv1.WalletOpRequest) (*userv1.EmptyResponse, error) {
	if err := biz.NewWalletLogic(s.svcCtx).SettleOrder(ctx, req.GetUserId(), req.GetAmount(), req.GetOrderId(), req.GetOrderNo()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}

func (s *UserServer) Notify(ctx context.Context, req *userv1.NotifyRequest) (*userv1.EmptyResponse, error) {
	if req.GetUserId() == 0 || req.GetTitle() == "" {
		return &userv1.EmptyResponse{}, nil
	}
	_, err := biz.NewUserLogic(s.svcCtx).CreateNotification(ctx, biz.NotifyCreateReq{
		UserID:   req.GetUserId(),
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		MsgType:  req.GetMsgType(),
		LinkType: req.GetLinkType(),
		LinkID:   req.GetLinkId(),
		Extra:    req.GetExtra(),
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}

func (s *UserServer) TaskEvent(ctx context.Context, req *userv1.TaskEventRequest) (*userv1.EmptyResponse, error) {
	if req.GetUserId() == 0 || req.GetTaskCode() == "" {
		return &userv1.EmptyResponse{}, nil
	}
	delta := int(req.GetDelta())
	if delta < 1 {
		delta = 1
	}
	if err := biz.NewTaskLogic(s.svcCtx).HandleEvent(ctx, biz.TaskEventReq{
		UserID:   req.GetUserId(),
		TaskCode: req.GetTaskCode(),
		Delta:    delta,
		RefType:  req.GetRefType(),
		RefID:    req.GetRefId(),
	}); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}

func StartZRPC(port int, svcCtx *svc.ServiceContext, logger *zap.Logger) *zrpc.RpcServer {
	c := zrpc.RpcServerConf{
		ListenOn: fmt.Sprintf("0.0.0.0:%d", port),
	}
	c.Mode = service.DevMode
	c.Log.Mode = "console"
	c.Log.Encoding = "plain"

	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		userv1.RegisterUserServiceServer(grpcServer, NewUserServer(svcCtx))
	})
	s.AddUnaryInterceptors(interceptor.Logging(logger))
	return s
}
