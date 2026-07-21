package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicGetArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicGetArticleLogic {
	return &PublicGetArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicGetArticleLogic) PublicGetArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 && in.Request != nil {
		if raw := in.Request.Header.Get(middleware.GatewayUserIDHeader); raw != "" {
			userID, _ = strconv.ParseUint(raw, 10, 64)
		}
	}
	data, err := clogic.NewArticleLogic(l.svcCtx).PublicDetail(ctx, id, userID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.AnyResp{Data: data}, nil
}
