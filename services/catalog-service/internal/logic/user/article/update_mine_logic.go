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

type UpdateMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMineLogic {
	return &UpdateMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMineLogic) UpdateMine(ctx context.Context, req *types.UserArticleUpdateBodyReq) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).UserUpdate(ctx, userID, id, req.ToContent()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
