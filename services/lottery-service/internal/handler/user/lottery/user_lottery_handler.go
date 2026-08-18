package lottery

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/lottery-service/internal/biz"
	"mymall/services/lottery-service/internal/svc"
	"mymall/services/lottery-service/internal/types"
)

func GetActivityHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID uint64
		if id, ok := middleware.GetUserID(r.Context()); ok {
			userID = id
		} else if s := r.Header.Get("X-User-Id"); s != "" {
			userID, _ = strconv.ParseUint(s, 10, 64)
		}
		vo, err := biz.NewLotteryLogic(svcCtx).GetCurrentActivity(r.Context(), userID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.DataResp{Data: vo})
	}
}

func DrawHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok || userID == 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
			return
		}
		vo, err := biz.NewLotteryLogic(svcCtx).Draw(r.Context(), userID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.DataResp{Data: vo})
	}
}

func ListRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok || userID == 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
			return
		}
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		list, total, err := biz.NewLotteryLogic(svcCtx).ListMyRecords(r.Context(), userID, req.Page, req.PageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
	}
}

func ClaimAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok || userID == 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
			return
		}
		var req types.ClaimAddressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		vo, err := biz.NewLotteryLogic(svcCtx).ClaimAddress(r.Context(), userID, req.Id, req.AddressID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.DataResp{Data: vo})
	}
}
