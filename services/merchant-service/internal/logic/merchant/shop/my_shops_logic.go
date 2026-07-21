package shop

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MyShopsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMyShopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MyShopsLogic {
	return &MyShopsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MyShopsLogic) MyShops(ctx context.Context) (resp *types.AnyResp, err error) {

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	shops, err := biz.NewMerchantLogic(l.svcCtx).MyShops(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: shops}, nil
}
