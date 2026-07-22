package banner

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBannerLogic {
	return &CreateBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateBannerLogic) CreateBanner(ctx context.Context, req *types.BannerSaveReq) (resp *types.BannerResp, err error) {
	b, err := clogic.NewArticleLogic(l.svcCtx).AdminCreateBanner(clogic.BannerSaveReq{
		Title: req.Title, ImageURL: req.ImageURL, LinkType: req.LinkType, LinkID: req.LinkID,
		Sort: req.Sort, Status: req.Status, StartAt: req.StartAt, EndAt: req.EndAt,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.BannerResp{Data: b}, nil
}
