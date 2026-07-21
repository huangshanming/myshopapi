package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteArticleCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteArticleCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteArticleCategoryLogic {
	return &AdminDeleteArticleCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteArticleCategoryLogic) AdminDeleteArticleCategory(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := clogic.NewArticleLogic(l.svcCtx).DeleteCategory(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
