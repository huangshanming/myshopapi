package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	shopopshandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
)

type AuthMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthMeLogic {
	return &AuthMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthMeLogic) AuthMe(w http.ResponseWriter, r *http.Request) {
	shopopshandler.NewShopOpsHandler(l.svcCtx).AuthMe(w, r)
}
