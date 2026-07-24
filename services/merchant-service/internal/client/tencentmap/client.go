package tencentmap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	key    string
	base   string
	client *http.Client
}

func New(key, baseURL string) *Client {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = "https://apis.map.qq.com"
	}
	return &Client{
		key:  strings.TrimSpace(key),
		base: base,
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (c *Client) Configured() bool {
	return c != nil && c.key != ""
}

type GeocoderResult struct {
	Province string  `json:"province"`
	City     string  `json:"city"`
	District string  `json:"district"`
	Address  string  `json:"address"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

func (c *Client) ReverseGeocode(lat, lng float64) (*GeocoderResult, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("未配置腾讯地图 Key")
	}
	u, _ := url.Parse(c.base + "/ws/geocoder/v1/")
	q := u.Query()
	q.Set("location", fmt.Sprintf("%f,%f", lat, lng))
	q.Set("key", c.key)
	q.Set("get_poi", "0")
	u.RawQuery = q.Encode()

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Address            string `json:"address"`
			FormattedAddresses struct {
				Recommend       string `json:"recommend"`
				Rough           string `json:"rough"`
				StandardAddress string `json:"standard_address"`
			} `json:"formatted_addresses"`
			AddressComponent struct {
				Province string `json:"province"`
				City     string `json:"city"`
				District string `json:"district"`
			} `json:"address_component"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if raw.Status != 0 {
		msg := raw.Message
		if msg == "" {
			msg = "逆地理编码失败"
		}
		return nil, fmt.Errorf("%s", msg)
	}
	city := raw.Result.AddressComponent.City
	if city == "" {
		city = raw.Result.AddressComponent.Province
	}
	addr := raw.Result.Address
	if fa := raw.Result.FormattedAddresses; fa.StandardAddress != "" {
		addr = fa.StandardAddress
	} else if fa.Recommend != "" {
		addr = fa.Recommend
	}
	return &GeocoderResult{
		Province: raw.Result.AddressComponent.Province,
		City:     city,
		District: raw.Result.AddressComponent.District,
		Address:  addr,
		Lat:      lat,
		Lng:      lng,
	}, nil
}
