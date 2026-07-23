package biz

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"mymall/services/user-service/internal/client/jutuike"
	"mymall/services/user-service/internal/svc"
)

type CpsLogic struct {
	svcCtx *svc.ServiceContext
}

func NewCpsLogic(svcCtx *svc.ServiceContext) *CpsLogic {
	return &CpsLogic{svcCtx: svcCtx}
}

type CpsActVO struct {
	ActID             uint64 `json:"act_id"`
	Title             string `json:"title"`
	Desc              string `json:"desc"`
	Img               string `json:"img"`
	Icon              string `json:"icon"`
	Poster            string `json:"poster"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	ActivityDate      string `json:"activity_date"`
	CommissionRateDes string `json:"commission_rate_des"`
	Introduce         string `json:"introduce"`
	SettlementTime    string `json:"settlement_time"`
}

type CpsConvertVO struct {
	H5        string         `json:"h5"`
	LongH5    string         `json:"long_h5"`
	ActName   string         `json:"act_name"`
	WeAppInfo map[string]any `json:"we_app_info,omitempty"`
}

func (l *CpsLogic) ListActs(ctx context.Context, page, pageSize int) ([]CpsActVO, int64, error) {
	cli := l.svcCtx.Jutuike
	if cli == nil || !cli.Configured() {
		return nil, 0, fmt.Errorf("未配置聚推客")
	}
	res, err := cli.ActList(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]CpsActVO, 0, len(res.List))
	for _, it := range res.List {
		out = append(out, toActVO(it))
	}
	return out, res.Total, nil
}

func (l *CpsLogic) Convert(ctx context.Context, userID, actID uint64) (*CpsConvertVO, error) {
	cli := l.svcCtx.Jutuike
	if cli == nil || !cli.Configured() {
		return nil, fmt.Errorf("未配置聚推客")
	}
	if userID == 0 {
		return nil, fmt.Errorf("未登录")
	}
	if actID == 0 {
		return nil, fmt.Errorf("活动ID无效")
	}
	sid := fmt.Sprintf("u%d", userID)
	if len(sid) > 30 {
		sid = sid[:30]
	}
	res, err := cli.ConvertAct(ctx, actID, sid)
	if err != nil {
		return nil, err
	}
	h5 := strings.TrimSpace(res.H5)
	if h5 == "" {
		h5 = strings.TrimSpace(res.ShortH5)
	}
	longH5 := strings.TrimSpace(res.LongH5)
	if h5 == "" && longH5 == "" {
		return nil, fmt.Errorf("该活动暂无可用推广链接")
	}
	return &CpsConvertVO{
		H5:        h5,
		LongH5:    longH5,
		ActName:   res.ActName,
		WeAppInfo: res.WeAppInfo,
	}, nil
}

func toActVO(it jutuike.ActItem) CpsActVO {
	id := parseUint(it.ActID)
	if id == 0 {
		id = parseUint(it.CouponID)
	}
	title := strings.TrimSpace(it.ActName)
	if title == "" {
		title = strings.TrimSpace(it.Title)
	}
	if title == "" {
		title = strings.TrimSpace(it.XcxShortTitle)
	}
	img := firstNonEmpty(it.Img, it.BgImages, it.Poster)
	poster := firstNonEmpty(it.Poster, it.BgImages, it.Img)
	activityDate := strings.TrimSpace(it.ActivityDate)
	if activityDate == "" {
		activityDate = joinDateRange(it.StartDate, it.EndDate)
	}
	commission := strings.TrimSpace(it.CommissionRateDes)
	if commission == "" {
		commission = strings.TrimSpace(it.SettlementTime)
	}
	return CpsActVO{
		ActID:             id,
		Title:             title,
		Desc:              strings.TrimSpace(it.Desc),
		Img:               img,
		Icon:              strings.TrimSpace(it.Icon),
		Poster:            poster,
		StartDate:         strings.TrimSpace(it.StartDate),
		EndDate:           strings.TrimSpace(it.EndDate),
		ActivityDate:      activityDate,
		CommissionRateDes: commission,
		Introduce:         strings.TrimSpace(it.Introduce),
		SettlementTime:    strings.TrimSpace(it.SettlementTime),
	}
}

type CpsGoodsVO struct {
	Platform       string `json:"platform"`
	ItemID         string `json:"item_id"`
	Title          string `json:"title"`
	Cover          string `json:"cover"`
	Price          string `json:"price"`
	CouponPrice    string `json:"coupon_price"`
	CommissionRate string `json:"commission_rate"`
	RawURL         string `json:"raw_url"`
}

type CpsGoodsConvertVO struct {
	H5     string `json:"h5"`
	LongH5 string `json:"long_h5"`
}

func (l *CpsLogic) ListGoods(ctx context.Context, platform, keyword string, page, pageSize int) ([]CpsGoodsVO, int64, error) {
	cli := l.svcCtx.Dingdanxia
	if cli == nil || !cli.Configured() {
		return nil, 0, fmt.Errorf("未配置订单侠")
	}
	platform = strings.ToLower(strings.TrimSpace(platform))
	res, err := cli.Search(ctx, platform, keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]CpsGoodsVO, 0, len(res.List))
	for _, it := range res.List {
		out = append(out, CpsGoodsVO{
			Platform:       it.Platform,
			ItemID:         it.ItemID,
			Title:          it.Title,
			Cover:          it.Cover,
			Price:          it.Price,
			CouponPrice:    it.CouponPrice,
			CommissionRate: it.CommissionRate,
			RawURL:         it.RawURL,
		})
	}
	return out, res.Total, nil
}

func (l *CpsLogic) ConvertGoods(ctx context.Context, userID uint64, platform, itemID, rawURL string) (*CpsGoodsConvertVO, error) {
	cli := l.svcCtx.Dingdanxia
	if cli == nil || !cli.Configured() {
		return nil, fmt.Errorf("未配置订单侠")
	}
	if userID == 0 {
		return nil, fmt.Errorf("未登录")
	}
	platform = strings.ToLower(strings.TrimSpace(platform))
	itemID = strings.TrimSpace(itemID)
	if itemID == "" && strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("商品ID无效")
	}
	sid := fmt.Sprintf("u%d", userID)
	if len(sid) > 30 {
		sid = sid[:30]
	}
	res, err := cli.Convert(ctx, platform, itemID, rawURL, sid)
	if err != nil {
		return nil, err
	}
	return &CpsGoodsConvertVO{H5: res.H5, LongH5: res.LongH5}, nil
}

func parseUint(n interface{ String() string }) uint64 {
	if n == nil {
		return 0
	}
	s := strings.TrimSpace(n.String())
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func joinDateRange(start, end string) string {
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	switch {
	case start != "" && end != "":
		return start + " 至 " + end
	case start != "":
		return start
	default:
		return end
	}
}
