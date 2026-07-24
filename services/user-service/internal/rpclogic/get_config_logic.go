package rpclogic

import (
	"context"
	"errors"
	"strings"

	userv1 "mymall/api/gen/user/v1"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigLogic {
	return &GetConfigLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetConfigLogic) GetConfig(in *userv1.GetConfigRequest) (*userv1.GetConfigResponse, error) {
	key := strings.TrimSpace(in.GetKey())
	if key == "" {
		return nil, status.Errorf(codes.InvalidArgument, "config key required")
	}
	val, err := biz.NewRBACLogic(l.svcCtx).GetConfigValue(l.ctx, key)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "config not found")
		}
		return nil, status.Errorf(codes.Internal, "config lookup failed")
	}
	return &userv1.GetConfigResponse{Key: key, Value: val}, nil
}
