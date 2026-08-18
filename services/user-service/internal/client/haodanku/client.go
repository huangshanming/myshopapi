package haodanku

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
	Pid     string // 淘宝推广位 PID
	TbName  string // 好单库已授权淘宝账号昵称
}

type Client struct {
	apiKey  string
	baseURL string
	pid     string
	tbName  string
	http    *http.Client
}

func NewClient(cfg Config) *Client {
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if base == "" {
		base = "https://v3.api.haodanku.com"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	return &Client{
		apiKey:  strings.TrimSpace(cfg.ApiKey),
		baseURL: base,
		pid:     strings.TrimSpace(cfg.Pid),
		tbName:  strings.TrimSpace(cfg.TbName),
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) Configured() bool {
	return c != nil && c.apiKey != ""
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
}

type SearchResult struct {
	Total     int64
	List      []GoodsItem
	NextMinID int64
}

type ConvertResult struct {
	H5     string
	LongH5 string
}

type ConvertOpts struct {
	ItemID string
	RawURL string
	Title  string
	Sid    string // 跟单参数 u{userID}
}

type apiEnvelope struct {
	Code  json.RawMessage `json:"code"`
	Msg   string          `json:"msg"`
	Data  json.RawMessage `json:"data"`
	MinID json.RawMessage `json:"min_id"`
}

func (c *Client) Search(ctx context.Context, platform, keyword string, minID, pageSize int) (*SearchResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置好单库")
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("请输入搜索关键词")
	}
	if minID < 1 {
		minID = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	switch platform {
	case "taobao":
		return c.searchTaobao(ctx, keyword, minID, pageSize)
	case "jd":
		return c.searchJD(ctx, keyword, minID, pageSize)
	case "pdd":
		return c.searchPDD(ctx, keyword, minID, pageSize)
	case "vip":
		return c.searchVIP(ctx, keyword, minID, pageSize)
	default:
		return nil, fmt.Errorf("不支持的平台")
	}
}

func (c *Client) Convert(ctx context.Context, platform string, opts ConvertOpts) (*ConvertResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置好单库")
	}
	opts.ItemID = strings.TrimSpace(opts.ItemID)
	opts.RawURL = strings.TrimSpace(opts.RawURL)
	opts.Title = strings.TrimSpace(opts.Title)
	opts.Sid = strings.TrimSpace(opts.Sid)
	if opts.ItemID == "" && opts.RawURL == "" {
		return nil, fmt.Errorf("商品ID无效")
	}
	switch platform {
	case "taobao":
		return c.convertTaobao(ctx, opts)
	case "jd":
		return c.convertJD(ctx, opts)
	case "pdd":
		return c.convertPDD(ctx, opts)
	case "vip":
		return c.convertVIP(ctx, opts)
	default:
		return nil, fmt.Errorf("不支持的平台")
	}
}

// ItemDetail fetches taobao item title when convert is missing title.
func (c *Client) ItemDetail(ctx context.Context, itemID string) (*GoodsItem, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置好单库")
	}
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return nil, fmt.Errorf("商品ID无效")
	}
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("itemid", itemID)
	env, err := c.get(ctx, "/item_detail", q)
	if err != nil {
		return nil, err
	}
	m := asMap(env.Data)
	if m == nil {
		return nil, fmt.Errorf("商品详情为空")
	}
	return &GoodsItem{
		Platform:       "taobao",
		ItemID:         firstStr(m, "itemid", "item_id"),
		Title:          firstStr(m, "itemtitle", "itemshorttitle", "title"),
		Cover:          firstStr(m, "itempic", "item_pic"),
		Price:          firstStr(m, "itemprice"),
		CouponPrice:    firstStr(m, "itemendprice"),
		CommissionRate: formatRate(firstStr(m, "tkrates")),
	}, nil
}

func (c *Client) searchTaobao(ctx context.Context, keyword string, minID, pageSize int) (*SearchResult, error) {
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("keyword", keyword)
	q.Set("min_id", strconv.Itoa(minID))
	q.Set("back", strconv.Itoa(pageSize))
	env, err := c.get(ctx, "/supersearch", q)
	if err != nil {
		return nil, err
	}
	return mapTBList(env, "taobao", minID, pageSize)
}

func (c *Client) searchJD(ctx context.Context, keyword string, minID, pageSize int) (*SearchResult, error) {
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("keyword", keyword)
	q.Set("min_id", strconv.Itoa(minID))
	q.Set("back", strconv.Itoa(pageSize))
	env, err := c.get(ctx, "/unify_jdgoods_search", q)
	if err != nil {
		return nil, err
	}
	rows := pickList(env.Data)
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "itemid", "item_id", "skuId", "sku_id")
		title := firstStr(row, "goodsname", "goodsName", "itemtitle", "title")
		cover := firstStr(row, "itempic", "item_pic", "pic")
		price := firstStr(row, "itemprice", "price")
		coupon := firstStr(row, "itemendprice", "coupon_price")
		rate := formatRate(firstStr(row, "jdrates", "commission_rate", "commissionRate"))
		raw := firstStr(row, "item_url", "materialUrl", "url")
		if id == "" && title == "" {
			continue
		}
		list = append(list, GoodsItem{
			Platform: "jd", ItemID: id, Title: title, Cover: cover,
			Price: price, CouponPrice: coupon, CommissionRate: rate, RawURL: raw,
		})
	}
	return &SearchResult{
		Total:     int64(len(list)),
		List:      list,
		NextMinID: nextMinID(env, minID, len(list), pageSize),
	}, nil
}

func (c *Client) searchPDD(ctx context.Context, keyword string, minID, pageSize int) (*SearchResult, error) {
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("keyword", keyword)
	q.Set("min_id", strconv.Itoa(minID))
	q.Set("limit", strconv.Itoa(pageSize))
	env, err := c.get(ctx, "/unify_pdd_goods_search", q)
	if err != nil {
		return nil, err
	}
	rows := pickList(env.Data)
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "goods_sign", "goodsSign", "itemid", "goods_id")
		title := firstStr(row, "goodsname", "goodsName", "itemtitle", "title")
		cover := firstStr(row, "itempic", "item_pic", "pdd_image")
		price := firstStr(row, "itemprice", "price")
		coupon := firstStr(row, "itemendprice")
		rate := formatRate(firstStr(row, "promotion_rate", "commission_rate", "tkrates"))
		if id == "" && title == "" {
			continue
		}
		list = append(list, GoodsItem{
			Platform: "pdd", ItemID: id, Title: title, Cover: cover,
			Price: price, CouponPrice: coupon, CommissionRate: rate,
		})
	}
	return &SearchResult{
		Total:     int64(len(list)),
		List:      list,
		NextMinID: nextMinID(env, minID, len(list), pageSize),
	}, nil
}

func (c *Client) searchVIP(ctx context.Context, keyword string, minID, pageSize int) (*SearchResult, error) {
	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("keyword", keyword)
	q.Set("min_id", strconv.Itoa(minID))
	q.Set("back", strconv.Itoa(pageSize))
	env, err := c.get(ctx, "/unify_vip_item_query", q)
	if err != nil {
		return nil, err
	}
	rows := pickList(env.Data)
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "goodsId", "goods_id", "itemid")
		title := firstStr(row, "goodsName", "goods_name", "shortTitle", "title")
		cover := firstStr(row, "goodsMainPicture", "itempic", "pic")
		price := firstStr(row, "marketPrice", "itemprice")
		coupon := firstStr(row, "vipPrice", "itemendprice")
		rate := formatRate(firstStr(row, "commissionRate", "commission_rate"))
		if id == "" && title == "" {
			continue
		}
		list = append(list, GoodsItem{
			Platform: "vip", ItemID: id, Title: title, Cover: cover,
			Price: price, CouponPrice: coupon, CommissionRate: rate,
		})
	}
	return &SearchResult{
		Total:     int64(len(list)),
		List:      list,
		NextMinID: nextMinID(env, minID, len(list), pageSize),
	}, nil
}

func mapTBList(env *apiEnvelope, platform string, curMinID, pageSize int) (*SearchResult, error) {
	rows := pickList(env.Data)
	list := make([]GoodsItem, 0, len(rows))
	for _, row := range rows {
		id := firstStr(row, "itemid", "item_id")
		title := firstStr(row, "itemtitle", "itemshorttitle", "title")
		cover := firstStr(row, "itempic", "item_pic")
		price := firstStr(row, "itemprice")
		coupon := firstStr(row, "itemendprice")
		rate := formatRate(firstStr(row, "tkrates", "commission_rate"))
		raw := firstStr(row, "item_url", "coupon_share_url", "url")
		if id == "" && title == "" {
			continue
		}
		list = append(list, GoodsItem{
			Platform: platform, ItemID: id, Title: title, Cover: cover,
			Price: price, CouponPrice: coupon, CommissionRate: rate, RawURL: raw,
		})
	}
	return &SearchResult{
		Total:     int64(len(list)),
		List:      list,
		NextMinID: nextMinID(env, curMinID, len(list), pageSize),
	}, nil
}

func (c *Client) convertTaobao(ctx context.Context, opts ConvertOpts) (*ConvertResult, error) {
	if c.pid == "" || c.tbName == "" {
		return nil, fmt.Errorf("未配置好单库淘宝推广位(Pid/TbName)")
	}
	title := opts.Title
	itemID := firstNonEmpty(opts.ItemID, opts.RawURL)
	if title == "" && opts.ItemID != "" {
		detail, err := c.ItemDetail(ctx, opts.ItemID)
		if err == nil && detail != nil {
			title = detail.Title
		}
	}
	if title == "" {
		return nil, fmt.Errorf("转链需要商品标题")
	}
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("itemid", itemID)
	params.Set("pid", c.pid)
	params.Set("tb_name", c.tbName)
	params.Set("get_taoword", "1")
	params.Set("title", title)
	if opts.Sid != "" {
		params.Set("relation_id", opts.Sid)
	}
	env, err := c.postForm(ctx, "/ratesurl", params)
	if err != nil {
		// 部分账号未开通渠道时 relation_id 会失败，去掉重试一次
		if opts.Sid != "" {
			params.Del("relation_id")
			env, err = c.postForm(ctx, "/ratesurl", params)
		}
		if err != nil {
			return nil, err
		}
	}
	m := asMap(env.Data)
	h5 := firstStr(m, "coupon_click_url", "item_url", "coupon_share_url", "short_url", "url")
	long := firstStr(m, "item_url", "coupon_click_url", "long_url")
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) convertJD(ctx context.Context, opts ConvertOpts) (*ConvertResult, error) {
	material := firstNonEmpty(opts.RawURL, opts.ItemID)
	if material != "" && !strings.Contains(material, "http") && opts.ItemID != "" {
		material = "https://item.jd.com/" + opts.ItemID + ".html"
	}
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("material_id", material)
	if opts.Sid != "" {
		params.Set("subUnionId", opts.Sid)
	}
	env, err := c.postForm(ctx, "/unify_jditems_link", params)
	if err != nil {
		return nil, err
	}
	m := asMap(env.Data)
	h5 := firstStr(m, "shortURL", "shortUrl", "clickURL", "clickUrl", "url")
	long := firstStr(m, "clickURL", "clickUrl", "longUrl")
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) convertPDD(ctx context.Context, opts ConvertOpts) (*ConvertResult, error) {
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("itemid", firstNonEmpty(opts.ItemID, opts.RawURL))
	if opts.Sid != "" {
		sid := opts.Sid
		if len(sid) > 20 {
			sid = sid[:20]
		}
		params.Set("channel", sid)
	}
	env, err := c.postForm(ctx, "/unify_pdditems_link", params)
	if err != nil {
		return nil, err
	}
	m := asMap(env.Data)
	h5 := firstStr(m, "mobile_short_url", "short_url", "mobile_url", "url")
	long := firstStr(m, "url", "mobile_url", "long_url")
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) convertVIP(ctx context.Context, opts ConvertOpts) (*ConvertResult, error) {
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("goods_id", firstNonEmpty(opts.ItemID, opts.RawURL))
	if opts.Sid != "" {
		sid := opts.Sid
		if len(sid) > 15 {
			sid = sid[:15]
		}
		params.Set("channel", sid)
	}
	env, err := c.postForm(ctx, "/unify_vip_item_convert", params)
	if err != nil {
		return nil, err
	}
	m := asMap(env.Data)
	h5 := firstStr(m, "url", "shortUrl", "longUrl", "deeplinkUrl")
	long := firstStr(m, "longUrl", "url")
	if h5 == "" && long == "" {
		return nil, fmt.Errorf("该商品暂无可用推广链接")
	}
	return &ConvertResult{H5: firstNonEmpty(h5, long), LongH5: long}, nil
}

func (c *Client) get(ctx context.Context, path string, q url.Values) (*apiEnvelope, error) {
	u := c.baseURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) postForm(ctx context.Context, path string, params url.Values) (*apiEnvelope, error) {
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req)
}

func (c *Client) do(req *http.Request) (*apiEnvelope, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求好单库失败: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("好单库响应无效: %w", err)
	}
	code := parseCode(env.Code)
	if code != 1 && code != 200 {
		msg := strings.TrimSpace(env.Msg)
		if msg == "" {
			msg = "好单库接口失败"
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return &env, nil
}

func nextMinID(env *apiEnvelope, curMinID, listLen, pageSize int) int64 {
	if n := parseInt64(env.MinID); n > 0 {
		// 好单库常返回「下一页」游标；若本页已空则结束
		if listLen == 0 {
			return 0
		}
		return n
	}
	if listLen == 0 || listLen < pageSize {
		return 0
	}
	return int64(curMinID + 1)
}

func parseCode(raw json.RawMessage) int {
	s := strings.Trim(string(raw), `"`)
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

func parseInt64(raw json.RawMessage) int64 {
	if len(raw) == 0 || string(raw) == "null" {
		return 0
	}
	s := strings.Trim(string(raw), `"`)
	n, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
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

func pickList(data json.RawMessage) []map[string]any {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var arr []map[string]any
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr
	}
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

func firstStr(m map[string]any, keys ...string) string {
	if m == nil {
		return ""
	}
	for _, k := range keys {
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
				// pdd_image 等可能是数组，取第一个
				if len(s) > 1 && s[0] == '[' {
					var ss []string
					if json.Unmarshal(b, &ss) == nil && len(ss) > 0 {
						return strings.TrimSpace(ss[0])
					}
				}
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
		return fmt.Sprintf("%.2f%%", n)
	}
	return s
}
