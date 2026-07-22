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

type ListBannersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListBannersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBannersLogic {
	return &ListBannersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListBannersLogic) ListBanners(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	list, err := clogic.NewArticleLogic(l.svcCtx).PublicBanners()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: int64(len(list))}, nil

}
