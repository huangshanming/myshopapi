package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pkgconfig "mymall/pkg/config"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	RpcListenOn string   `json:",default=0.0.0.0:9092"`
	Etcd        EtcdConf `json:",optional"`

	MySQL     MySQLConf
	// AppTelemetry is our OTel init config. Named to avoid clash with rest.RestConf / ServiceConf.Telemetry.
	AppTelemetry TelemetryConf
	TencentMap   TencentMapConf `json:",optional"`
}

type TencentMapConf struct {
	Key     string `json:",optional"`
	BaseURL string `json:",default=https://apis.map.qq.com"`
}

type EtcdConf struct {
	Hosts []string `json:",optional"`
}

type MySQLConf struct {
	Host            string `json:",default=127.0.0.1"`
	Port            int    `json:",default=3306"`
	Username        string `json:",default=homestead"`
	Password        string
	Dbname          string `json:",default=mymall"`
	Charset         string `json:",default=utf8mb4"`
	MaxOpenConns    int    `json:",default=50"`
	MaxIdleConns    int    `json:",default=10"`
	ConnMaxLifetime int    `json:",default=3600"`
	ConnMaxIdleTime int    `json:",default=3600"`
}

type TelemetryConf struct {
	Enabled  bool   `json:",default=false"`
	Endpoint string `json:",optional"`
	Service  string `json:",default=merchant-service"`
}

func (c MySQLConf) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Dbname, c.Charset,
	)
}

func (c MySQLConf) ToPkg() pkgconfig.MySQLConfig {
	return pkgconfig.MySQLConfig{
		Username: c.Username, Password: c.Password, Host: c.Host, Port: c.Port,
		Dbname: c.Dbname, Charset: c.Charset,
		MaxOpenConns: c.MaxOpenConns, MaxIdleConns: c.MaxIdleConns,
		ConnMaxLifetime: c.ConnMaxLifetime, ConnMaxIdleTime: c.ConnMaxIdleTime,
	}
}

func (c TelemetryConf) ToPkg() pkgconfig.TelemetryConfig {
	return pkgconfig.TelemetryConfig{Enabled: c.Enabled, Endpoint: c.Endpoint, Service: c.Service}
}


func (c Config) GRPCPort() int {
	parts := strings.Split(c.RpcListenOn, ":")
	if len(parts) == 0 {
		return 9092
	}
	p, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil || p == 0 {
		return 9092
	}
	return p
}

func (c *Config) OverlayFromEnv() {
	if v := strings.TrimSpace(os.Getenv("MYMALL_SERVER_HTTP_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.Port = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_SERVER_GRPC_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.RpcListenOn = fmt.Sprintf("0.0.0.0:%d", n)
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_SERVER_MODE")); v != "" {
		switch strings.ToLower(v) {
		case "release", "prod", "production", "pro":
			c.Mode = "pro"
		default:
			c.Mode = "dev"
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_MYSQL_HOST")); v != "" {
		c.MySQL.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_MYSQL_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.MySQL.Port = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_MYSQL_USERNAME")); v != "" {
		c.MySQL.Username = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_MYSQL_PASSWORD")); v != "" {
		c.MySQL.Password = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_MYSQL_DBNAME")); v != "" {
		c.MySQL.Dbname = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_MYSQL_CHARSET")); v != "" {
		c.MySQL.Charset = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_ETCD_HOSTS")); v != "" {
		c.Etcd.Hosts = splitHosts(v)
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TELEMETRY_ENABLED")); v != "" {
		c.AppTelemetry.Enabled = v == "1" || strings.EqualFold(v, "true")
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TELEMETRY_ENDPOINT")); v != "" {
		c.AppTelemetry.Endpoint = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TELEMETRY_SERVICE")); v != "" {
		c.AppTelemetry.Service = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TENCENT_MAP_KEY")); v != "" {
		c.TencentMap.Key = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TENCENT_MAP_BASE_URL")); v != "" {
		c.TencentMap.BaseURL = v
	}
}

func splitHosts(v string) []string {
	parts := strings.FieldsFunc(v, func(r rune) bool {
		return r == ',' || r == ' ' || r == ';'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
