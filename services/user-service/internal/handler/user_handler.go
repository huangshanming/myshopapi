package handler

import (
	"encoding/json"
	"net/http"

	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type UserHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.UserLogic
}

func NewUserHandler(svcCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{
		svcCtx: svcCtx,
		logic:  logic.NewUserLogic(svcCtx),
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	token, user, err := h.logic.LoginWithShop(req.Mobile, req.Password, req.ShopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	response.Success(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"mobile":   user.Mobile,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
			"status":   user.Status,
		},
	}, "登录成功")
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	user, err := h.logic.Register(req.Mobile, req.Password)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, user, "注册成功")
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	var userID uint64
	if id, ok := middleware.GetUserID(r.Context()); ok {
		userID = id
	} else if claims, ok := jwt.ClaimsFromContext(r.Context()); ok {
		userID = claims.UserID
	} else {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	user, err := h.logic.GetProfile(userID)
	if err != nil {
		response.Error(w, "用户不存在", http.StatusNotFound)
		return
	}
	response.Success(w, user, "查询成功")
}
