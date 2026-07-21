package article

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminPurgeArticleRecycleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminPurgeArticleRecycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminPurgeArticleRecycleLogic {
	return &AdminPurgeArticleRecycleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminPurgeArticleRecycleLogic) AdminPurgeArticleRecycle(ctx context.Context) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{}

	var body struct {
		ID uint64 `json:"id"`
	}
	_ = appinput.BindBody(in, &body)
	if body.ID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "缺少 id")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).PermanentDelete(ctx, body.ID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
