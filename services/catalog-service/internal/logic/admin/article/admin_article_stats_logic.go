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

type AdminArticleStatsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminArticleStatsLogic {
	return &AdminArticleStatsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminArticleStatsLogic) AdminArticleStats(ctx context.Context) (resp *types.AnyResp, err error) {

	data, err := clogic.NewArticleLogic(l.svcCtx).Stats(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
