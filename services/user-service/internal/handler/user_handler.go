package handler

import (
	"net/http"

	"mymall/pkg/jwt"
	"mymall/pkg/response"
	"mymall/services/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type loginReq struct {
	Mobile   string `json:"mobile" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6"`
}

type registerReq struct {
	Mobile   string `json:"mobile" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	token, user, err := h.svc.Login(req.Mobile, req.Password)
	if err != nil {
		response.Error(c, err.Error(), http.StatusUnauthorized)
		return
	}
	response.Success(c, gin.H{"token": token, "user": user}, "登录成功")
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	user, err := h.svc.Register(req.Mobile, req.Password)
	if err != nil {
		response.Error(c, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(c, user, "注册成功")
}

func (h *UserHandler) Profile(c *gin.Context) {
	claims, ok := c.Get("jwt_claims")
	if !ok {
		response.Error(c, "未授权", http.StatusUnauthorized)
		return
	}
	userClaims := claims.(*jwt.Claims)
	user, err := h.svc.GetProfile(userClaims.UserID)
	if err != nil {
		response.Error(c, "用户不存在", http.StatusNotFound)
		return
	}
	response.Success(c, user, "查询成功")
}
