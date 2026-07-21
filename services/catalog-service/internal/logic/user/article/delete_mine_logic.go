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

type DeleteMineLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMineLogic {
	return &DeleteMineLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeleteMineLogic) DeleteMine(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "文章ID无效")
	}
	if err := clogic.NewArticleLogic(l.svcCtx).UserDelete(ctx, userID, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
