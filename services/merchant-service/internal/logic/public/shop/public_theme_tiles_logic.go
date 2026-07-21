package shop

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicThemeTilesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicThemeTilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicThemeTilesLogic {
	return &PublicThemeTilesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicThemeTilesLogic) PublicThemeTiles(ctx context.Context) (resp *types.AnyResp, err error) {

	list, err := biz.NewMerchantLogic(l.svcCtx).ListThemeTiles()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: map[string]interface{}{"list": list}}, nil
}
