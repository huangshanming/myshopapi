package article

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDetailMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailMineLogic {
	return &DetailMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DetailMineLogic) DetailMine(ctx context.Context, req *types.IdPathReq) (resp *types.ArticleResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id := req.Id
	data, err := clogic.NewArticleLogic(l.svcCtx).UserGetMine(ctx, userID, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.ArticleResp{Data: data}, nil
}
