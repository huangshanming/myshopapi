package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	ctypes "mymall/services/catalog-service/internal/content/types"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateArticleLogic {
	return &AdminUpdateArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateArticleLogic) AdminUpdateArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

	uid, _ := middleware.GetUserID(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body ctypes.ArticleSaveReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).AdminUpdate(ctx, id, uid, body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
