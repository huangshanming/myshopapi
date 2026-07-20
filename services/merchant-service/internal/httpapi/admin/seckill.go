package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *SeckillHandler) AdminGetSeckillRule(w http.ResponseWriter, r *http.Request) {
	rule, err := h.logic.GetSeckillRule()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, rule)
}

func (h *SeckillHandler) AdminUpdateSeckillRule(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillRuleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	rule, err := h.logic.UpdateSeckillRule(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, rule)
}

func (h *SeckillHandler) AdminListSeckillSessions(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListSeckillSessions(p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *SeckillHandler) AdminListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	sid, _ := strconv.ParseUint(r.URL.Query().Get("session_id"), 10, 64)
	list, total, err := h.logic.ListAdminSeckillEntries(sid, p, ps)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}
