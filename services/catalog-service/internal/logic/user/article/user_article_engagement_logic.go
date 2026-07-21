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
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewArticleHandler(l.svcCtx).Status(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
