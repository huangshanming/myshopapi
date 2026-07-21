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

type AdminGetArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetArticleLogic {
	return &AdminGetArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGetArticleLogic) AdminGetArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	data, err := clogic.NewArticleLogic(l.svcCtx).Detail(ctx, id, 0)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
