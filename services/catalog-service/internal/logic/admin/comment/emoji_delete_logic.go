package comment

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type EmojiDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmojiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiDeleteLogic {
	return &EmojiDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmojiDeleteLogic) EmojiDelete(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).EmojiDelete(w, r)
}
