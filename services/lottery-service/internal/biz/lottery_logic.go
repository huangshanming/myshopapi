package biz

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mymall/pkg/cache"
	"mymall/services/lottery-service/internal/client/userhttp"
	"mymall/services/lottery-service/internal/model"
	"mymall/services/lottery-service/internal/repository"
	"mymall/services/lottery-service/internal/svc"
	"mymall/services/lottery-service/internal/uploadpath"
)

type LotteryLogic struct {
	svcCtx *svc.ServiceContext
}

func NewLotteryLogic(svcCtx *svc.ServiceContext) *LotteryLogic {
	return &LotteryLogic{svcCtx: svcCtx}
}

type PrizeVO struct {
	ID           uint64 `json:"id"`
	Slot         int    `json:"slot"`
	Name         string `json:"name"`
	CoverURL     string `json:"cover_url"`
	PrizeType    string `json:"prize_type"`
	PointsAmount int    `json:"points_amount"`
	Weight       int    `json:"weight,omitempty"`
	Stock        int    `json:"stock,omitempty"`
	StockStrict  int    `json:"stock_strict"`
}

type ActivityVO struct {
	ID             uint64    `json:"id"`
	Title          string    `json:"title"`
	Status         int       `json:"status"`
	CostPoints     int       `json:"cost_points"`
	DailyLimit     int       `json:"daily_limit"`
	StartAt        string    `json:"start_at"`
	EndAt          string    `json:"end_at"`
	Prizes         []PrizeVO `json:"prizes"`
	TodayUsed      int64     `json:"today_used"`
	TodayRemaining int64     `json:"today_remaining"`
}

type DrawResultVO struct {
	RecordID      uint64 `json:"record_id"`
	Slot          int    `json:"slot"`
	PrizeID       uint64 `json:"prize_id"`
	PrizeName     string `json:"prize_name"`
	PrizeType     string `json:"prize_type"`
	PointsAmount  int    `json:"points_amount"`
	CostPoints    int    `json:"cost_points"`
	FulfillStatus string `json:"fulfill_status"`
	CoverURL      string `json:"cover_url,omitempty"`
}

type RecordVO struct {
	ID              uint64 `json:"id"`
	ActivityID      uint64 `json:"activity_id"`
	UserID          uint64 `json:"user_id,omitempty"`
	Slot            int    `json:"slot"`
	PrizeName       string `json:"prize_name"`
	PrizeType       string `json:"prize_type"`
	PointsAmount    int    `json:"points_amount"`
	CostPoints      int    `json:"cost_points"`
	FulfillStatus   string `json:"fulfill_status"`
	ReceiverName    string `json:"receiver_name,omitempty"`
	ReceiverPhone   string `json:"receiver_phone,omitempty"`
	ReceiverAddress string `json:"receiver_address,omitempty"`
	ShipCompany     string `json:"ship_company,omitempty"`
	ShipNo          string `json:"ship_no,omitempty"`
	ShippedAt       string `json:"shipped_at,omitempty"`
	CreatedAt       string `json:"created_at"`
}

type ActivitySaveReq struct {
	Title      string
	Status     int
	CostPoints int
	DailyLimit int
	StartAt    string
	EndAt      string
}

type PrizeSaveItem struct {
	Slot         int
	Name         string
	CoverURL     string
	PrizeType    string
	PointsAmount int
	Weight       int
	Stock        int
	StockStrict  int
}

func (l *LotteryLogic) GetCurrentActivity(ctx context.Context, userID uint64) (*ActivityVO, error) {
	a, err := l.svcCtx.Repo.GetActiveActivity(ctx)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("暂无进行中的抽奖活动")
	}
	return l.toActivityVO(ctx, a, userID, false)
}

func (l *LotteryLogic) toActivityVO(ctx context.Context, a *model.LotteryActivity, userID uint64, admin bool) (*ActivityVO, error) {
	prizes, err := l.svcCtx.Repo.ListPrizes(ctx, a.ID)
	if err != nil {
		return nil, err
	}
	vo := &ActivityVO{
		ID: a.ID, Title: a.Title, Status: a.Status,
		CostPoints: a.CostPoints, DailyLimit: a.DailyLimit,
		StartAt: repository.FormatTimePtr(a.StartAt),
		EndAt:   repository.FormatTimePtr(a.EndAt),
		Prizes:  make([]PrizeVO, 0, len(prizes)),
	}
	for _, p := range prizes {
		item := PrizeVO{
			ID: p.ID, Slot: p.Slot, Name: p.Name, CoverURL: p.CoverURL,
			PrizeType: p.PrizeType, PointsAmount: p.PointsAmount,
		}
		if admin {
			item.Weight = p.Weight
			item.Stock = p.Stock
			item.StockStrict = p.StockStrict
		}
		vo.Prizes = append(vo.Prizes, item)
	}
	if userID > 0 {
		used, err := l.svcCtx.Repo.CountUserDrawsToday(ctx, userID, a.ID)
		if err != nil {
			return nil, err
		}
		vo.TodayUsed = used
		if a.DailyLimit > 0 {
			left := int64(a.DailyLimit) - used
			if left < 0 {
				left = 0
			}
			vo.TodayRemaining = left
		} else {
			vo.TodayRemaining = -1
		}
	}
	return vo, nil
}

func (l *LotteryLogic) Draw(ctx context.Context, userID uint64) (*DrawResultVO, error) {
	if userID == 0 {
		return nil, errors.New("未登录")
	}
	a, err := l.svcCtx.Repo.GetActiveActivity(ctx)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("暂无进行中的抽奖活动")
	}
	if a.DailyLimit > 0 {
		used, err := l.svcCtx.Repo.CountUserDrawsToday(ctx, userID, a.ID)
		if err != nil {
			return nil, err
		}
		if used >= int64(a.DailyLimit) {
			return nil, errors.New("今日抽奖次数已用完")
		}
	}
	prizes, err := l.svcCtx.Repo.ListPrizes(ctx, a.ID)
	if err != nil {
		return nil, err
	}
	if len(prizes) != 9 {
		return nil, errors.New("活动奖品未配置完整（需要9格）")
	}

	recID, err := l.svcCtx.Repo.InsertPendingRecord(ctx, &model.LotteryDrawRecord{
		UserID: userID, ActivityID: a.ID, CostPoints: a.CostPoints,
	})
	if err != nil {
		return nil, err
	}

	cost := a.CostPoints
	if cost > 0 {
		if err := l.svcCtx.UserHTTP.Deduct(ctx, ledger(userID, int64(cost), "lottery_cost", "lottery_draw", recID, "九宫格抽奖消耗")); err != nil {
			_ = l.svcCtx.Repo.MarkRecordFailed(ctx, recID)
			return nil, err
		}
	}

	idx, err := weightedPick(prizes)
	if err != nil {
		l.rollbackCost(ctx, userID, cost, recID)
		_ = l.svcCtx.Repo.MarkRecordFailed(ctx, recID)
		return nil, err
	}

	// 有限库存：先 Redis 预扣（inventory-sync 预热 + Canal binlog 对齐），再落 MySQL。
	chosen := prizes[idx]
	redisDeducted := false
	if chosen.Stock >= 0 {
		if err := cache.LotteryStockDeduct(ctx, l.svcCtx.Redis, chosen.ID, 1); err != nil {
			l.rollbackCost(ctx, userID, cost, recID)
			_ = l.svcCtx.Repo.MarkRecordFailed(ctx, recID)
			if errors.Is(err, cache.ErrStockInsufficient) {
				return nil, errors.New("奖品已抽完")
			}
			return nil, fmt.Errorf("库存预扣失败: %w", err)
		}
		redisDeducted = true
	}

	prize, err := l.svcCtx.Repo.FinalizeDraw(ctx, recID, prizes, idx)
	if err != nil {
		if redisDeducted {
			_ = cache.LotteryStockRestore(ctx, l.svcCtx.Redis, chosen.ID, 1)
		}
		l.rollbackCost(ctx, userID, cost, recID)
		_ = l.svcCtx.Repo.MarkRecordFailed(ctx, recID)
		return nil, err
	}

	if prize.PrizeType == model.PrizeTypePoints && prize.PointsAmount > 0 {
		if err := l.svcCtx.UserHTTP.Add(ctx, ledger(userID, int64(prize.PointsAmount), "lottery_reward", "lottery_draw", recID, "九宫格抽奖奖励："+prize.Name)); err != nil {
			l.rollbackCost(ctx, userID, cost, recID)
			return nil, fmt.Errorf("发奖失败: %w", err)
		}
	}

	fulfill := model.FulfillNone
	if prize.PrizeType == model.PrizeTypePhysical {
		fulfill = model.FulfillNeedAddress
	}
	return &DrawResultVO{
		RecordID: recID, Slot: prize.Slot, PrizeID: prize.ID,
		PrizeName: prize.Name, PrizeType: prize.PrizeType,
		PointsAmount: prize.PointsAmount, CostPoints: cost,
		FulfillStatus: fulfill, CoverURL: prize.CoverURL,
	}, nil
}

func (l *LotteryLogic) rollbackCost(ctx context.Context, userID uint64, cost int, recID uint64) {
	if cost <= 0 {
		return
	}
	_ = l.svcCtx.UserHTTP.Refund(ctx, ledger(userID, int64(cost), "lottery_cost_refund", "lottery_draw", recID, "九宫格抽奖退还"))
}

func ledger(userID uint64, points int64, changeType, refType string, refID uint64, reason string) userhttp.LedgerReq {
	return userhttp.LedgerReq{
		UserID: userID, Points: points, Reason: reason,
		RefNo: strconv.FormatUint(refID, 10), ChangeType: changeType, RefType: refType,
	}
}

func weightedPick(prizes []model.LotteryPrize) (int, error) {
	total := 0
	eligible := make([]int, 0, len(prizes))
	weights := make([]int, 0, len(prizes))
	for i, p := range prizes {
		if p.Weight <= 0 || p.Stock == 0 {
			continue
		}
		eligible = append(eligible, i)
		weights = append(weights, p.Weight)
		total += p.Weight
	}
	if total <= 0 {
		return 0, errors.New("暂无可用奖品")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(total)))
	if err != nil {
		return 0, err
	}
	cur := int(n.Int64())
	for i, w := range weights {
		if cur < w {
			return eligible[i], nil
		}
		cur -= w
	}
	return eligible[len(eligible)-1], nil
}

func toRecordVO(r model.LotteryDrawRecord, withUser bool) RecordVO {
	vo := RecordVO{
		ID: r.ID, ActivityID: r.ActivityID, Slot: r.Slot,
		PrizeName: r.PrizeName, PrizeType: r.PrizeType,
		PointsAmount: r.PointsAmount, CostPoints: r.CostPoints,
		FulfillStatus:   r.FulfillStatus,
		ReceiverName:    r.ReceiverName,
		ReceiverPhone:   r.ReceiverPhone,
		ReceiverAddress: r.ReceiverAddress,
		ShipCompany:     r.ShipCompany,
		ShipNo:          r.ShipNo,
		CreatedAt:       r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if withUser {
		vo.UserID = r.UserID
	}
	if r.ShippedAt.Valid {
		vo.ShippedAt = r.ShippedAt.Time.Format("2006-01-02 15:04:05")
	}
	if vo.FulfillStatus == "" {
		vo.FulfillStatus = model.FulfillNone
	}
	return vo
}

func (l *LotteryLogic) ListMyRecords(ctx context.Context, userID uint64, page, pageSize int) ([]RecordVO, int64, error) {
	list, total, err := l.svcCtx.Repo.ListUserRecords(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]RecordVO, 0, len(list))
	for _, r := range list {
		out = append(out, toRecordVO(r, false))
	}
	return out, total, nil
}

func (l *LotteryLogic) ClaimAddress(ctx context.Context, userID, recordID, addressID uint64) (*RecordVO, error) {
	if userID == 0 || recordID == 0 || addressID == 0 {
		return nil, errors.New("参数无效")
	}
	rec, err := l.svcCtx.Repo.GetRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if rec == nil || rec.UserID != userID {
		return nil, errors.New("记录不存在")
	}
	if rec.Status != model.RecordStatusDone || rec.PrizeType != model.PrizeTypePhysical {
		return nil, errors.New("该记录不是实物奖")
	}
	if rec.FulfillStatus != model.FulfillNeedAddress {
		return nil, errors.New("地址已填写或无需填写")
	}
	addr, err := l.svcCtx.UserHTTP.GetAddress(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(addr.ReceiverName)
	phone := strings.TrimSpace(addr.ReceiverPhone)
	full := strings.TrimSpace(addr.FullAddress())
	if name == "" || phone == "" || full == "" {
		return nil, errors.New("收货地址信息不完整")
	}
	if err := l.svcCtx.Repo.ClaimAddress(ctx, recordID, userID, addressID, name, phone, full); err != nil {
		return nil, err
	}
	updated, err := l.svcCtx.Repo.GetRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}
	vo := toRecordVO(*updated, false)
	return &vo, nil
}

func (l *LotteryLogic) AdminListActivities(ctx context.Context, page, pageSize int) ([]ActivityVO, int64, error) {
	list, total, err := l.svcCtx.Repo.ListActivities(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]ActivityVO, 0, len(list))
	for i := range list {
		vo, err := l.toActivityVO(ctx, &list[i], 0, true)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *vo)
	}
	return out, total, nil
}

func (l *LotteryLogic) AdminGetActivity(ctx context.Context, id uint64) (*ActivityVO, error) {
	a, err := l.svcCtx.Repo.GetActivity(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("活动不存在")
	}
	return l.toActivityVO(ctx, a, 0, true)
}

func (l *LotteryLogic) AdminCreateActivity(ctx context.Context, req ActivitySaveReq) (*ActivityVO, error) {
	a, err := buildActivity(0, req)
	if err != nil {
		return nil, err
	}
	id, err := l.svcCtx.Repo.CreateActivity(ctx, a)
	if err != nil {
		return nil, err
	}
	a.ID = id
	return l.toActivityVO(ctx, a, 0, true)
}

func (l *LotteryLogic) AdminUpdateActivity(ctx context.Context, id uint64, req ActivitySaveReq) (*ActivityVO, error) {
	exist, err := l.svcCtx.Repo.GetActivity(ctx, id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, errors.New("活动不存在")
	}
	a, err := buildActivity(id, req)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.Repo.UpdateActivity(ctx, a); err != nil {
		return nil, err
	}
	return l.toActivityVO(ctx, a, 0, true)
}

func buildActivity(id uint64, req ActivitySaveReq) (*model.LotteryActivity, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("活动标题不能为空")
	}
	if req.CostPoints < 0 {
		return nil, errors.New("消耗积分不能为负")
	}
	if req.DailyLimit < 0 {
		return nil, errors.New("每日次数不能为负")
	}
	start, err := repository.ParseTimePtr(req.StartAt)
	if err != nil {
		return nil, err
	}
	end, err := repository.ParseTimePtr(req.EndAt)
	if err != nil {
		return nil, err
	}
	status := req.Status
	if status < 0 || status > 2 {
		status = model.ActivityStatusDraft
	}
	return &model.LotteryActivity{
		ID: id, Title: title, Status: status,
		CostPoints: req.CostPoints, DailyLimit: req.DailyLimit,
		StartAt: start, EndAt: end,
	}, nil
}

func (l *LotteryLogic) AdminSavePrizes(ctx context.Context, activityID uint64, items []PrizeSaveItem) (*ActivityVO, error) {
	a, err := l.svcCtx.Repo.GetActivity(ctx, activityID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("活动不存在")
	}
	if len(items) != 9 {
		return nil, errors.New("必须配置恰好9个奖品位")
	}
	seen := map[int]bool{}
	prizes := make([]model.LotteryPrize, 0, 9)
	for _, it := range items {
		if it.Slot < 0 || it.Slot > 8 {
			return nil, errors.New("slot 必须是 0-8")
		}
		if seen[it.Slot] {
			return nil, errors.New("slot 不能重复")
		}
		seen[it.Slot] = true
		pt := strings.TrimSpace(it.PrizeType)
		if pt == "" {
			pt = model.PrizeTypeThanks
		}
		if pt != model.PrizeTypePoints && pt != model.PrizeTypeThanks && pt != model.PrizeTypePhysical {
			return nil, errors.New("prize_type 仅支持 points/thanks/physical")
		}
		name := strings.TrimSpace(it.Name)
		if name == "" {
			return nil, errors.New("奖品名称不能为空")
		}
		if it.Weight < 0 {
			return nil, errors.New("权重不能为负")
		}
		pointsAmt := it.PointsAmount
		if pt == model.PrizeTypePoints {
			if pointsAmt <= 0 {
				return nil, errors.New("积分奖品数量须大于0")
			}
		} else {
			pointsAmt = 0
		}
		strict := 0
		if it.StockStrict != 0 {
			strict = 1
		}
		prizes = append(prizes, model.LotteryPrize{
			ActivityID: activityID, Slot: it.Slot, Name: name, CoverURL: strings.TrimSpace(it.CoverURL),
			PrizeType: pt, PointsAmount: pointsAmt, Weight: it.Weight, Stock: it.Stock, StockStrict: strict,
		})
	}
	for i := 0; i < 9; i++ {
		if !seen[i] {
			return nil, fmt.Errorf("缺少 slot=%d", i)
		}
	}
	if err := l.svcCtx.Repo.ReplacePrizes(ctx, activityID, prizes); err != nil {
		return nil, err
	}
	return l.toActivityVO(ctx, a, 0, true)
}

func (l *LotteryLogic) AdminListRecords(ctx context.Context, activityID uint64, prizeType string, page, pageSize int) ([]RecordVO, int64, error) {
	pt := strings.TrimSpace(prizeType)
	if pt != "" && pt != model.PrizeTypePoints && pt != model.PrizeTypeThanks && pt != model.PrizeTypePhysical {
		return nil, 0, errors.New("prize_type 仅支持 points/thanks/physical")
	}
	list, total, err := l.svcCtx.Repo.ListRecordsAdmin(ctx, activityID, pt, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]RecordVO, 0, len(list))
	for _, r := range list {
		out = append(out, toRecordVO(r, true))
	}
	return out, total, nil
}

func (l *LotteryLogic) AdminListOrders(ctx context.Context, fulfillStatus string, page, pageSize int) ([]RecordVO, int64, error) {
	list, total, err := l.svcCtx.Repo.ListFulfillmentOrders(ctx, fulfillStatus, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]RecordVO, 0, len(list))
	for _, r := range list {
		out = append(out, toRecordVO(r, true))
	}
	return out, total, nil
}

func (l *LotteryLogic) AdminShipOrder(ctx context.Context, id uint64, company, shipNo string) error {
	company = strings.TrimSpace(company)
	shipNo = strings.TrimSpace(shipNo)
	if shipNo == "" {
		return errors.New("物流单号不能为空")
	}
	rec, err := l.svcCtx.Repo.GetRecord(ctx, id)
	if err != nil {
		return err
	}
	if rec == nil {
		return errors.New("订单不存在")
	}
	if err := l.svcCtx.Repo.ShipRecord(ctx, id, company, shipNo); err != nil {
		return err
	}
	content := fmt.Sprintf("您的抽奖奖品「%s」已发货", rec.PrizeName)
	if company != "" {
		content += fmt.Sprintf("，物流：%s %s", company, shipNo)
	} else {
		content += fmt.Sprintf("，运单号：%s", shipNo)
	}
	_ = l.svcCtx.UserHTTP.Notify(ctx, userhttp.NotifyReq{
		UserID:  rec.UserID,
		Title:   "抽奖奖品已发货",
		Content: content,
		Type:    "lottery",
	})
	return nil
}

func (l *LotteryLogic) SaveUpload(filename string, data []byte) (string, error) {
	if len(data) > 5*1024*1024 {
		return "", errors.New("文件不能超过5MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return "", errors.New("仅支持图片")
	}
	dir := uploadpath.Abs("lottery")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return "/uploads/lottery/" + name, nil
}
