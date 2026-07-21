package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminRestoreArticleRecycleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminRestoreArticleRecycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRestoreArticleRecycleLogic {
	return &AdminRestoreArticleRecycleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminRestoreArticleRecycleLogic) AdminRestoreArticleRecycle(ctx context.Context, req *types.ArticleIdListReq) (resp *types.AnyResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "缺少 id")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).Restore(ctx, req.Id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
