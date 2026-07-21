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

type UpdateBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBannerLogic {
	return &UpdateBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateBannerLogic) UpdateBanner(ctx context.Context, req *types.BannerUpdateBodyReq) (resp *types.AnyResp, err error) {
	if err := clogic.NewArticleLogic(l.svcCtx).AdminUpdateBanner(req.Id, clogic.BannerSaveReq{Title: req.Title, ImageURL: req.ImageURL, LinkType: req.LinkType, LinkID: req.LinkID, Sort: req.Sort, Status: req.Status, StartAt: req.StartAt, EndAt: req.EndAt}); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
