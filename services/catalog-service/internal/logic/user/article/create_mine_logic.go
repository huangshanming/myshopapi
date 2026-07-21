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

type CreateMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMineLogic {
	return &CreateMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateMineLogic) CreateMine(ctx context.Context, req *types.UserArticleCreateReq) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	a, err := clogic.NewArticleLogic(l.svcCtx).UserCreate(ctx, userID, req.ToContent())
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: a}, nil
}
