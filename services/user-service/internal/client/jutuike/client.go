package jutuike

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ApiKey  string
	BaseURL string
	Timeout time.Duration
}

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func NewClient(cfg Config) *Client {
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if base == "" {
		base = "http://api.jutuike.com"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	return &Client{
		apiKey:  strings.TrimSpace(cfg.ApiKey),
		baseURL: base,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) Configured() bool {
	return c != nil && c.apiKey != ""
}

type ActItem struct {
	ActID               json.Number `json:"act_id"`
	ActName             string      `json:"act_name"`
	Desc                string      `json:"desc"`
	Img                 string      `json:"img"`
	Icon                string      `json:"icon"`
	Poster              string      `json:"poster"`
	StartDate           string      `json:"start_date"`
	EndDate             string      `json:"end_date"`
	Introduce           string      `json:"introduce"`
	AttributionExplain  string      `json:"attribution_explain"`
	Note                string      `json:"note"`
	SettlementTime      string      `json:"settlement_time"`
	CommissionRateDes   string      `json:"commission_rate_des"`
	ActivityDate        string      `json:"activity_date"`
	Title               string      `json:"title"`
	XcxShortTitle       string      `json:"xcx_short_title"`
	BgImages            string      `json:"bg_images"`
	CouponID            json.Number `json:"coupon_id"`
}

type ActListResult struct {
	Total int64
	List  []ActItem
}

type ConvertResult struct {
	H5        string         `json:"h5"`
	LongH5    string         `json:"long_h5"`
	ShortH5   string         `json:"short_h5"`
	Deeplink  string         `json:"deeplink"`
	ActName   string         `json:"act_name"`
	WeAppInfo map[string]any `json:"we_app_info"`
}

type apiEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (c *Client) ActList(ctx context.Context, page, pageSize int) (*ActListResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置聚推客")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	body, err := c.get(ctx, "/union/act_list", q)
	if err != nil {
		return nil, err
	}
	return parseActList(body)
}

func (c *Client) ConvertAct(ctx context.Context, actID uint64, sid string) (*ConvertResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置聚推客")
	}
	if actID == 0 {
		return nil, fmt.Errorf("活动ID无效")
	}
	sid = strings.TrimSpace(sid)
	if sid == "" {
		return nil, fmt.Errorf("跟单参数无效")
	}
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("sid", sid)
	q.Set("act_id", strconv.FormatUint(actID, 10))
	body, err := c.get(ctx, "/union/act", q)
	if err != nil {
		return nil, err
	}
	var out ConvertResult
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("解析转链结果失败: %w", err)
	}
	return &out, nil
}

func (c *Client) get(ctx context.Context, path string, q url.Values) (json.RawMessage, error) {
	u := c.baseURL + path + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求聚推客失败: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("聚推客响应无效: %w", err)
	}
	if env.Code != 1 {
		msg := strings.TrimSpace(env.Msg)
		if msg == "" {
			msg = "聚推客接口失败"
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return env.Data, nil
}

func parseActList(data json.RawMessage) (*ActListResult, error) {
	data = json.RawMessage(strings.TrimSpace(string(data)))
	if len(data) == 0 || string(data) == "null" || string(data) == `""` {
		return &ActListResult{List: []ActItem{}}, nil
	}
	if data[0] == '[' {
		var list []ActItem
		if err := json.Unmarshal(data, &list); err != nil {
			return nil, fmt.Errorf("解析活动列表失败: %w", err)
		}
		if list == nil {
			list = []ActItem{}
		}
		return &ActListResult{Total: int64(len(list)), List: list}, nil
	}
	var obj struct {
		Total       json.Number     `json:"total"`
		List        json.RawMessage `json:"list"`
		Data        json.RawMessage `json:"data"`
		CurrentPage json.Number     `json:"current_page"`
		PerPage     json.Number     `json:"per_page"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("解析活动列表失败: %w", err)
	}
	rawList := obj.List
	if len(rawList) == 0 {
		rawList = obj.Data
	}
	list := []ActItem{}
	if len(rawList) > 0 && string(rawList) != "null" {
		if err := json.Unmarshal(rawList, &list); err != nil {
			return nil, fmt.Errorf("解析活动列表失败: %w", err)
		}
		if list == nil {
			list = []ActItem{}
		}
	}
	total, _ := obj.Total.Int64()
	if total == 0 {
		total = int64(len(list))
	}
	return &ActListResult{Total: total, List: list}, nil
}
