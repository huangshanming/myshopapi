package handler

import (
	"net/http"

	"mymall/pkg/apidoc/dto"
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

// Login 用户登录
// @Summary      用户登录
// @Description  手机号 + 密码登录，返回 JWT Token
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginReq  true  "登录参数"
// @Success      200   {object}  apidoc.Response{data=dto.LoginResp}  "登录成功"
// @Router       /api/v1/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", http.StatusBadRequest)
		return
	}
	token, user, err := h.svc.LoginWithShop(req.Mobile, req.Password, req.ShopID)
	if err != nil {
		response.Error(c, err.Error(), http.StatusUnauthorized)
		return
	}
	response.Success(c, gin.H{"token": token, "user": user}, "登录成功")
}

// Register 用户注册
// @Summary      用户注册
// @Description  手机号注册新用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterReq  true  "注册参数"
// @Success      200   {object}  apidoc.Response{data=dto.UserInfo}  "注册成功"
// @Router       /api/v1/user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
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

// Profile 获取个人资料
// @Summary      获取个人资料
// @Description  需要 JWT 鉴权
// @Tags         用户
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  apidoc.Response{data=dto.UserProfileResp}  "查询成功"
// @Router       /api/v1/user/profile [get]
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
