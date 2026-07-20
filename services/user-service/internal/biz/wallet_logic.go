package biz

import (
	"context"
	"errors"

	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/svc"
)

type WalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletLogic {
	return &WalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WalletLogic) GetWallet(userID uint64) (*model.UserWallet, error) {
	if userID == 0 {
		return nil, errors.New("用户无效")
	}
	return l.svcCtx.Repo.GetWallet(userID)
}

func (l *WalletLogic) AdjustWallet(userID uint64, field string, amount float64, remark string, operatorID uint64) (*model.UserWallet, error) {
	if userID == 0 {
		return nil, errors.New("用户无效")
	}
	if amount == 0 {
		return nil, errors.New("调账金额不能为 0")
	}
	var op *uint64
	if operatorID > 0 {
		op = &operatorID
	}
	return l.svcCtx.Repo.AdjustWallet(userID, field, amount, remark, op)
}

func (l *WalletLogic) ListWalletLogs(userID uint64, page, pageSize int) ([]model.UserWalletLog, int64, error) {
	if userID == 0 {
		return nil, 0, errors.New("用户无效")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return l.svcCtx.Repo.ListWalletLogs(userID, page, pageSize)
}

func (l *WalletLogic) FreezeForOrder(userID uint64, amount float64, orderID uint64, orderNo string) error {
	if userID == 0 || orderID == 0 {
		return errors.New("参数无效")
	}
	return l.svcCtx.Repo.FreezeForOrder(userID, amount, orderID, orderNo)
}

func (l *WalletLogic) UnfreezeOrder(userID uint64, amount float64, orderID uint64, orderNo string) error {
	if userID == 0 || orderID == 0 {
		return errors.New("参数无效")
	}
	return l.svcCtx.Repo.UnfreezeOrder(userID, amount, orderID, orderNo)
}

func (l *WalletLogic) SettleOrder(userID uint64, amount float64, orderID uint64, orderNo string) error {
	if userID == 0 || orderID == 0 {
		return errors.New("参数无效")
	}
	return l.svcCtx.Repo.SettleOrder(userID, amount, orderID, orderNo)
}
