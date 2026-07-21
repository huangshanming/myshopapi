package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"
	"strconv"
	"time"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListArticlesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListArticlesLogic {
	return &AdminListArticlesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListArticlesLogic) AdminListArticles(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := req.Page, req.PageSize
	f := repository.ArticleListFilter{
		Title:       "" /* was query:title */,
		AuditStatus: "" /* was query:audit_status */,
		Status:      "" /* was query:status */,
		Page:        page, PageSize: pageSize,
		Recycle: false,
	}
	if s := "" /* was query:shop_id */; s != "" {
		shopID, _ := strconv.ParseUint(s, 10, 64)
		f.ShopID = shopID
		f.FilterShop = true
	}
	if s := "" /* was query:has_schedule */; s == "1" {
		v := true
		f.HasSchedule = &v
	} else if s == "0" {
		v := false
		f.HasSchedule = &v
	}
	if s := "" /* was query:created_from */; s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := "" /* was query:created_to */; s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	data, err := clogic.NewArticleLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
