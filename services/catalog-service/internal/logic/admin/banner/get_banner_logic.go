package banner

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBannerLogic {
	return &GetBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetBannerLogic) GetBanner(ctx context.Context, req *types.IdPathReq) (resp *types.BannerResp, err error) {
	id := req.Id
	b, err := clogic.NewArticleLogic(l.svcCtx).AdminGetBanner(id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.BannerResp{Data: b}, nil
}
