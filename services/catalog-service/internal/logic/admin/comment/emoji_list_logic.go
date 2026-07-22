package comment

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmojiListLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewEmojiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiListLogic {
	return &EmojiListLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *EmojiListLogic) EmojiList(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := req.Page, req.PageSize
	data, err := clogic.NewArticleLogic(l.svcCtx).ListEmojisAdmin(ctx, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.FromPaged(data), nil
}
