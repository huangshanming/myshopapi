package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hpublic "mymall/services/catalog-service/internal/content/app/public"
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

func (l *CreateMineLogic) CreateMine(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewArticleHandler(l.svcCtx).CreateMine(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
