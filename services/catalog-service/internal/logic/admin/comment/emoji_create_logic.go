package comment

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type EmojiCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmojiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiCreateLogic {
	return &EmojiCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmojiCreateLogic) EmojiCreate(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).EmojiCreate(w, r)
}
