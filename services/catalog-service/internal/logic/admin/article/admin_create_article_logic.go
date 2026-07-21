package article

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	ctypes "mymall/services/catalog-service/internal/content/types"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateArticleLogic {
	return &AdminCreateArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateArticleLogic) AdminCreateArticle(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	uid, _ := middleware.GetUserID(ctx)
	var body ctypes.ArticleSaveReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	a, err := clogic.NewArticleLogic(l.svcCtx).AdminCreate(ctx, uid, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: a}, nil
}
