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

type ListEmojisLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmojisLogic {
	return &ListEmojisLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListEmojisLogic) ListEmojis(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	list, err := clogic.NewArticleLogic(l.svcCtx).ListEmojisPublic(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: int64(len(list))}, nil

}
