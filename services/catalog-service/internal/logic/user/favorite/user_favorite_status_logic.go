package favorite

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	huser "mymall/services/catalog-service/internal/product/app/user"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFavoriteStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserFavoriteStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFavoriteStatusLogic {
	return &UserFavoriteStatusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserFavoriteStatusLogic) UserFavoriteStatus(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := huser.NewFavoriteHandler(l.svcCtx).Status(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
