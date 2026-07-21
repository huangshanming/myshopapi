package banner

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/content/logic"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBannerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBannerLogic {
	return &CreateBannerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateBannerLogic) CreateBanner(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body logic.BannerSaveReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	b, err := clogic.NewArticleLogic(l.svcCtx).AdminCreateBanner(body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: b}, nil
}
