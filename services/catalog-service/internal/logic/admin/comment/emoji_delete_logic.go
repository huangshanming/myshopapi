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

type EmojiDeleteLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewEmojiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiDeleteLogic {
	return &EmojiDeleteLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *EmojiDeleteLogic) EmojiDelete(ctx context.Context, req *types.IdPathReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if err := clogic.NewArticleLogic(l.svcCtx).DeleteEmoji(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
