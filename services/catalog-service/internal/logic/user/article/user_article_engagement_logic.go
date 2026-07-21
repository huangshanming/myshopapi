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

type UserArticleEngagementLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserArticleEngagementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserArticleEngagementLogic {
	return &UserArticleEngagementLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserArticleEngagementLogic) UserArticleEngagement(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id := req.Id
	liked, favorited := clogic.NewArticleLogic(l.svcCtx).EngagementStatus(ctx, userID, id)
	return &types.AnyResp{Data: map[string]bool{"liked": liked, "favorited": favorited}}, nil

}
