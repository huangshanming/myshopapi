package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"
	"time"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListArticleRecycleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListArticleRecycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListArticleRecycleLogic {
	return &AdminListArticleRecycleLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListArticleRecycleLogic) AdminListArticleRecycle(ctx context.Context, req *types.AdminArticleListReq) (resp *types.PageListResp, err error) {
	f := repository.ArticleListFilter{
		Title: req.Title, AuditStatus: req.AuditStatus, Status: req.Status,
		Page: req.Page, PageSize: req.PageSize, Recycle: true,
	}
	if req.ShopId > 0 {
		f.ShopID = req.ShopId
		f.FilterShop = true
	}
	if req.HasSchedule == "1" {
		v := true
		f.HasSchedule = &v
	} else if req.HasSchedule == "0" {
		v := false
		f.HasSchedule = &v
	}
	if req.CreatedFrom != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedFrom, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if req.CreatedTo != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedTo, time.Local); err == nil {
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
