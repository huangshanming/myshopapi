package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type GenerateUserTokenLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGenerateUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateUserTokenLogic {
	return &GenerateUserTokenLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GenerateUserTokenLogic) GenerateUserToken(ctx context.Context, req *types.IdPathReq) (resp *types.UserTokenResp, err error) {
	data, err := biz.NewRBACLogic(l.svcCtx).GenerateUserToken(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	out := &types.UserTokenResp{}
	if v, ok := data["token"].(string); ok {
		out.Token = v
	}
	if v, ok := data["user_id"].(uint64); ok {
		out.UserId = v
	}
	if v, ok := data["mobile"].(string); ok {
		out.Mobile = v
	}
	if v, ok := data["nickname"].(string); ok {
		out.Nickname = v
	}
	if v, ok := data["role"].(string); ok {
		out.Role = v
	}
	if v, ok := data["shop_id"].(uint64); ok {
		out.ShopId = v
	}
	switch v := data["expire_hours"].(type) {
	case int:
		out.ExpireHours = v
	case int64:
		out.ExpireHours = int(v)
	}
	return out, nil
}
