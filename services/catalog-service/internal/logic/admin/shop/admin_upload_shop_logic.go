package shop

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	hadmin "mymall/services/catalog-service/internal/product/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUploadShopLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUploadShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUploadShopLogic {
	return &AdminUploadShopLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUploadShopLogic) AdminUploadShop(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewShopUploadHandler().Upload(ctx, appinput.CallInput{Request: r})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
