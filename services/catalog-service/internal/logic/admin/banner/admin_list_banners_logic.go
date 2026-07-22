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

type AdminListBannersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListBannersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListBannersLogic {
	return &AdminListBannersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListBannersLogic) AdminListBanners(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := int(req.Page), int(req.PageSize)
	data, err := clogic.NewArticleLogic(l.svcCtx).AdminListBanners(page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.FromPaged(data), nil
}
