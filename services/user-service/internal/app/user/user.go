package user

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/types"
	"net/http"
	"strconv"
)

func (h *UserHandler) ListNotifications(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	page, pageSize := in.Page()
	list, total, err := h.logic.ListMyNotifications(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *UserHandler) UnreadNotificationCount(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	n, err := h.logic.UnreadCount(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"count": n}, nil
}

func (h *UserHandler) MarkNotificationRead(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.MarkRead(ctx, userID, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *UserHandler) MarkAllNotificationsRead(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	if err := h.logic.MarkAllRead(ctx, userID); err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func (h *UserHandler) Login(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.LoginReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	token, user, err := h.logic.LoginWithShop(ctx, req.Mobile, req.Password, req.ShopID)
	if err != nil {
		return nil, xerr.New(http.StatusUnauthorized, err.Error())
	}
	return map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"mobile":   user.Mobile,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
			"status":   user.Status,
		},
	}, nil
}

func (h *UserHandler) Register(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.RegisterReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	user, err := h.logic.Register(ctx, req.Mobile, req.Password)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return user, nil
}

func (h *UserHandler) Profile(ctx context.Context, in appinput.CallInput) (any, error) {
	var userID uint64
	if id, ok := middleware.GetUserID(ctx); ok {
		userID = id
	} else if claims, ok := jwt.ClaimsFromContext(ctx); ok {
		userID = claims.UserID
	} else {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	user, err := h.logic.GetProfile(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "用户不存在")
	}
	return user, nil
}
