package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr")

type UserHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.UserLogic
}

func NewUserHandler(svcCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{
		svcCtx: svcCtx,
		logic:  logic.NewUserLogic(context.Background(), svcCtx),
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	token, user, err := h.logic.LoginWithShop(req.Mobile, req.Password, req.ShopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"mobile":   user.Mobile,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
			"status":   user.Status,
		},
	})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	user, err := h.logic.Register(req.Mobile, req.Password)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, user)
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	var userID uint64
	if id, ok := middleware.GetUserID(r.Context()); ok {
		userID = id
	} else if claims, ok := jwt.ClaimsFromContext(r.Context()); ok {
		userID = claims.UserID
	} else {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未授权"))
		return
	}
	user, err := h.logic.GetProfile(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusNotFound, "用户不存在"))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, user)
}
