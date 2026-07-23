package dingdanxia

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
	PddPid  string
}

type Client struct {
	apiKey  string
	baseURL string
	pddPid  string
	http    *http.Client
}

func NewClient(cfg Config) *Client {
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if base == "" {
		base = "http://api.tbk.dingdanxia.com"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	return &Client{
		apiKey:  strings.TrimSpace(cfg.ApiKey),
		baseURL: base,
		pddPid:  strings.TrimSpace(cfg.PddPid),
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) Configured() bool {
	return c != nil && c.apiKey != ""
}

func (c *Client) PddPid() string {
	if c == nil {
		return ""
	}
	return c.pddPid
}

type GoodsItem struct {
	Platform       string
	ItemID         string
	Title          string
	Cover          string
	Price          string
	CouponPrice    string
	CommissionRate string
	RawURL         string
	Extra          map[string]string
}

type SearchResult struct {
	Total int64
	List  []GoodsItem
}

type ConvertResult struct {
	H5     string
	LongH5 string
}

type apiEnvelope struct {
	Code json.RawMessage `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (c *Client) Search(ctx context.Context, platform, keyword string, page, pageSize int) (*SearchResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置订单侠")
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("请输入搜索关键词")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	switch platform {
	case "taobao":
		return c.searchTaobao(ctx, keyword, page, pageSize)
	case "jd":
		return c.searchJD(ctx, keyword, page, pageSize)
	case "pdd":
		return c.searchPDD(ctx, keyword, page, pageSize)
	case "vip":
		return c.searchVIP(ctx, keyword, page, pageSize)
	default:
		return nil, fmt.Errorf("不支持的平台")
	}
}

func (c *Client) Convert(ctx context.Context, platform, itemID, rawURL, sid string) (*ConvertResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置订单侠")
	}
	sid = strings.TrimSpace(sid)
	if sid == "" {
		return nil, fmt.Errorf("跟单参数无效")
	}
	switch platform {
	case "taobao":
		return c.convertTaobao(ctx, itemID, rawURL, sid)
	case "jd":
		return c.convertJD(ctx, itemID, rawURL, sid)
	case "pdd":
		return c.convertPDD(ctx, itemID, sid)
	case "vip":
		return c.convertVIP(ctx, itemID, sid)
	default:
		return nil, fmt.Errorf("不支持的平台")
	}
}

func (c *Client) searchTaobao(ctx context.Context, keyword string, page, pageSize int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("q", keyword)
	params.Set("has_coupon", "false")
	params.Set("page_no", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	data, err := c.post(ctx, "/tbk/super_search", params)
	if err != nil {
		return nil, err
	}
	rows := asObjectSlice(data)
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "item_id", "num_iid", "itemId")
		title := firstStr(row, "title", "short_title")
		cover := firstStr(row, "pict_url", "white_image")
		price := firstStr(row, "zk_final_price", "reserve_price")
		coupon := firstStr(row, "coupon_amount", "coupon")
		rate := formatRate(firstStr(row, "commission_rate"))
		raw := firstStr(row, "coupon_share_url", "url", "item_url")
		if id == "" && title == "" {
			continue
		}
		list = append(list, GoodsItem{
			Platform: "taobao", ItemID: id, Title: title, Cover: cover,
			Price: price, CouponPrice: couponPrice(price, coupon), CommissionRate: rate, RawURL: raw,
		})
	}
	total := int64(len(list))
	if t := firstStr(asMap(data), "total_results", "total"); t != "" {
		if n, err := strconv.ParseInt(t, 10, 64); err == nil && n > 0 {
			total = n
		}
	}
	return &SearchResult{Total: total, List: list}, nil
}

func (c *Client) convertTaobao(ctx context.Context, itemID, rawURL, sid string) (*ConvertResult, error) {
	params := url.Values{}
	params.Set("id", firstNonEmpty(itemID, rawURL))
	params.Set("relation_id", sid)
	data, err := c.post(ctx, "/tbk/id_privilege", params)
	if err != nil {
		// fallback universal convert
		params2 := url.Values{}
		params2.Set("material_url", firstNonEmpty(rawURL, itemID))
		params2.Set("relation_id", sid)
		data, err = c.post(ctx, "/tbk/wn_convert", params2)
		if err != nil {
			return nil, err
		}
	}
	m := asMap(data)
	h5 := firstStr(m, "coupon_click_url", "item_url", "coupon_share_url", "short_url", "url", "click_url")
	long := firstStr(m, "max_commission_rate_url", "long_url", "item_url")
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) searchJD(ctx context.Context, keyword string, page, pageSize int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("pageIndex", strconv.Itoa(page))
	params.Set("pageSize", strconv.Itoa(pageSize))
	data, err := c.post(ctx, "/jd/query_goods", params)
	if err != nil {
		return nil, err
	}
	rows := pickList(data, "list", "data", "goods", "result")
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "skuId", "sku_id", "itemId", "item_id", "wareId")
		title := firstStr(row, "skuName", "sku_name", "goodsName", "title", "wareName")
		cover := firstStr(row, "imageUrl", "image_url", "picUrl", "pictUrl", "img")
		price := firstStr(row, "price", "wlPrice", "lowestPrice", "priceInfo.price")
		if price == "" {
			if pi, ok := row["priceInfo"].(map[string]any); ok {
				price = firstStr(pi, "lowestCouponPrice", "lowestPrice", "price")
			}
		}
		rate := ""
		if ci, ok := row["commissionInfo"].(map[string]any); ok {
			rate = firstStr(ci, "commissionShare", "commission")
			if rate != "" && !strings.Contains(rate, "%") {
				rate += "%"
			}
		}
		raw := firstStr(row, "materialUrl", "material_url", "itemUrl", "url")
		couponPrice := ""
		if pi, ok := row["priceInfo"].(map[string]any); ok {
			couponPrice = firstStr(pi, "lowestCouponPrice")
		}
		if id == "" && title == "" {
			continue
		}
		list = append(list, GoodsItem{
			Platform: "jd", ItemID: id, Title: title, Cover: cover,
			Price: price, CouponPrice: couponPrice, CommissionRate: rate, RawURL: raw,
		})
	}
	total := int64(len(list))
	if m := asMap(data); m != nil {
		if t := firstStr(m, "total", "totalCount", "totalNum"); t != "" {
			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				total = n
			}
		}
	}
	return &SearchResult{Total: total, List: list}, nil
}

func (c *Client) convertJD(ctx context.Context, itemID, rawURL, sid string) (*ConvertResult, error) {
	material := firstNonEmpty(rawURL, itemID)
	if material != "" && !strings.Contains(material, "http") && itemID != "" {
		material = "https://item.jd.com/" + itemID + ".html"
	}
	params := url.Values{}
	params.Set("materialId", material)
	params.Set("subUnionId", sid)
	data, err := c.post(ctx, "/jd/jy_privilege2", params)
	if err != nil {
		return nil, err
	}
	m := asMap(data)
	h5 := firstStr(m, "shortURL", "shortUrl", "clickURL", "clickUrl", "url")
	long := firstStr(m, "longUrl", "clickURL", "clickUrl")
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) searchPDD(ctx context.Context, keyword string, page, pageSize int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	if c.pddPid != "" {
		params.Set("pid", c.pddPid)
	}
	data, err := c.post(ctx, "/pdd/goods_search", params)
	if err != nil {
		return nil, err
	}
	rows := pickList(data, "goods_list", "list", "data")
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "goods_sign", "goodsSign", "goods_id", "goodsId")
		title := firstStr(row, "goods_name", "goodsName", "title")
		cover := firstStr(row, "goods_thumbnail_url", "goods_image_url", "thumb_url", "img")
		price := fenToYuan(firstStr(row, "min_group_price", "min_normal_price", "price"))
		coupon := fenToYuan(firstStr(row, "coupon_discount", "coupon_price"))
		rate := firstStr(row, "promotion_rate", "commission_rate")
		if rate != "" {
			if n, err := strconv.ParseFloat(rate, 64); err == nil && n > 100 {
				rate = fmt.Sprintf("%.2f%%", n/10)
			} else if !strings.Contains(rate, "%") {
				rate += "%"
			}
		}
		if id == "" && title == "" {
			continue
		}
		extra := map[string]string{}
		if gs := firstStr(row, "goods_sign", "goodsSign"); gs != "" {
			extra["goods_sign"] = gs
		}
		list = append(list, GoodsItem{
			Platform: "pdd", ItemID: id, Title: title, Cover: cover,
			Price: price, CouponPrice: couponPrice(price, coupon), CommissionRate: rate, Extra: extra,
		})
	}
	total := int64(len(list))
	if m := asMap(data); m != nil {
		if t := firstStr(m, "total_count", "total", "totalCount"); t != "" {
			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				total = n
			}
		}
	}
	return &SearchResult{Total: total, List: list}, nil
}

func (c *Client) convertPDD(ctx context.Context, itemID, sid string) (*ConvertResult, error) {
	if c.pddPid == "" {
		return nil, fmt.Errorf("未配置拼多多推广位PID")
	}
	params := url.Values{}
	params.Set("goods_sign", itemID)
	params.Set("p_id", c.pddPid)
	params.Set("custom_parameters", sid)
	data, err := c.post(ctx, "/pdd/pdd_convert", params)
	if err != nil {
		params2 := url.Values{}
		params2.Set("goods_sign_list", itemID)
		params2.Set("p_id_list", c.pddPid)
		params2.Set("custom_parameters", sid)
		data, err = c.post(ctx, "/pdd/goods_zs_unit_url_gen", params2)
		if err != nil {
			return nil, err
		}
	}
	m := asMap(data)
	h5 := firstStr(m, "mobile_short_url", "mobile_url", "short_url", "url", "mobile_url_list")
	long := firstStr(m, "url", "mobile_url", "long_url")
	if h5 == "" {
		// nested list style
		if arr := asObjectSlice(data); len(arr) > 0 {
			h5 = firstStr(arr[0], "mobile_short_url", "mobile_url", "short_url", "url")
			long = firstStr(arr[0], "url", "mobile_url")
		}
	}
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) searchVIP(ctx context.Context, keyword string, page, pageSize int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("page", strconv.Itoa(page))
	params.Set("pageSize", strconv.Itoa(pageSize))
	data, err := c.post(ctx, "/vip/search", params)
	if err != nil {
		return nil, err
	}
	rows := pickList(data, "goodsInfoList", "list", "data", "goodsList")
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "goodsId", "goods_id", "productId", "item_id")
		title := firstStr(row, "goodsName", "goods_name", "productName", "title")
		cover := firstStr(row, "goodsThumbUrl", "goodsMainPicture", "img", "picture")
		price := firstStr(row, "vipPrice", "price", "marketPrice")
		rate := firstStr(row, "commissionRate", "commission", "commRate")
		if rate != "" && !strings.Contains(rate, "%") {
			rate += "%"
		}
		if id == "" && title == "" {
			continue
		}
		list = append(list, GoodsItem{
			Platform: "vip", ItemID: id, Title: title, Cover: cover,
			Price: price, CommissionRate: rate,
		})
	}
	total := int64(len(list))
	if m := asMap(data); m != nil {
		if t := firstStr(m, "total", "totalCount"); t != "" {
			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				total = n
			}
		}
	}
	return &SearchResult{Total: total, List: list}, nil
}

func (c *Client) convertVIP(ctx context.Context, itemID, sid string) (*ConvertResult, error) {
	params := url.Values{}
	params.Set("goodsId", itemID)
	params.Set("chanTag", sid)
	data, err := c.post(ctx, "/vip/id_privilege", params)
	if err != nil {
		return nil, err
	}
	m := asMap(data)
	h5 := firstStr(m, "url", "shortUrl", "ulUrl", "longUrl")
	long := firstStr(m, "longUrl", "noEvokeLongUrl", "url")
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) post(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", c.apiKey)
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求订单侠失败: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("订单侠响应无效: %w", err)
	}
	code := parseCode(env.Code)
	if code != 200 && code != 1 {
		msg := strings.TrimSpace(env.Msg)
		if msg == "" {
			msg = "订单侠接口失败"
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return env.Data, nil
}

func parseCode(raw json.RawMessage) int {
	s := strings.Trim(string(raw), `"`)
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

func asMap(data json.RawMessage) map[string]any {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err == nil {
		return m
	}
	return nil
}

func asObjectSlice(data json.RawMessage) []map[string]any {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var arr []map[string]any
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr
	}
	// sometimes data is object wrapping list
	m := asMap(data)
	if m == nil {
		return nil
	}
	for _, k := range []string{"list", "data", "goods_list", "goodsInfoList", "result"} {
		if v, ok := m[k]; ok {
			b, _ := json.Marshal(v)
			var nested []map[string]any
			if err := json.Unmarshal(b, &nested); err == nil {
				return nested
			}
		}
	}
	return nil
}

func pickList(data json.RawMessage, keys ...string) []map[string]any {
	if rows := asObjectSlice(data); len(rows) > 0 && data[0] == '[' {
		return rows
	}
	m := asMap(data)
	if m == nil {
		return asObjectSlice(data)
	}
	for _, k := range keys {
		if v, ok := m[k]; ok {
			b, _ := json.Marshal(v)
			var nested []map[string]any
			if err := json.Unmarshal(b, &nested); err == nil {
				return nested
			}
		}
	}
	return asObjectSlice(data)
}

func firstStr(m map[string]any, keys ...string) string {
	if m == nil {
		return ""
	}
	for _, k := range keys {
		if strings.Contains(k, ".") {
			parts := strings.SplitN(k, ".", 2)
			if nested, ok := m[parts[0]].(map[string]any); ok {
				if s := firstStr(nested, parts[1]); s != "" {
					return s
				}
			}
			continue
		}
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			if s := strings.TrimSpace(t); s != "" {
				return s
			}
		case float64:
			if t == float64(int64(t)) {
				return strconv.FormatInt(int64(t), 10)
			}
			return strconv.FormatFloat(t, 'f', -1, 64)
		case json.Number:
			return t.String()
		case bool:
			return strconv.FormatBool(t)
		default:
			b, _ := json.Marshal(t)
			s := strings.Trim(string(b), `"`)
			if s != "" && s != "null" {
				return s
			}
		}
	}
	return ""
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func formatRate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "%") {
		return s
	}
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		if n > 100 {
			return fmt.Sprintf("%.2f%%", n/100)
		}
		return fmt.Sprintf("%.2f%%", n)
	}
	return s
}

func fenToYuan(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	if n >= 100 {
		return fmt.Sprintf("%.2f", n/100)
	}
	return fmt.Sprintf("%.2f", n)
}

func couponPrice(price, coupon string) string {
	p, err1 := strconv.ParseFloat(strings.TrimSpace(price), 64)
	c, err2 := strconv.ParseFloat(strings.TrimSpace(coupon), 64)
	if err1 != nil || err2 != nil || c <= 0 {
		return ""
	}
	v := p - c
	if v < 0 {
		v = 0
	}
	return fmt.Sprintf("%.2f", v)
}
