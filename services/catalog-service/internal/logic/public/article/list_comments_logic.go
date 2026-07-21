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

type ListCommentsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *ListCommentsLogic) ListComments(ctx context.Context, req *types.IdPageReq) (resp *types.PageListResp, err error) {
	data, err := clogic.NewArticleLogic(l.svcCtx).PublicListComments(ctx, req.Id, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
