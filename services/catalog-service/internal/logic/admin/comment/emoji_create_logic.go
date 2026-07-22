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

type EmojiCreateLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewEmojiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiCreateLogic {
	return &EmojiCreateLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *EmojiCreateLogic) EmojiCreate(ctx context.Context, req *types.EmojiSaveReq) (resp *types.EmojiResp, err error) {
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	e, err := clogic.NewArticleLogic(l.svcCtx).CreateEmoji(ctx, req.Name, req.ImageURL, req.Sort, status)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmojiResp{Data: e}, nil
}
