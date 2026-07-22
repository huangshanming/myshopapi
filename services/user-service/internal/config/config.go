package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pkgconfig "mymall/pkg/config"

	"github.com/zeromicro/go-zero/rest"
)

// Config is the go-zero style service config (RestConf + business fields).
type Config struct {
	rest.RestConf

	// RpcListenOn e.g. 0.0.0.0:9090
	RpcListenOn string `json:",default=0.0.0.0:9090"`

	MySQL     MySQLConf
	Auth      AuthConf
	Telemetry TelemetryConf
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

type AuthConf struct {
	AccessSecret string
	AccessExpire int64  `json:",default=86400"` // seconds
	ConsumerKey  string `json:",default=mymall-user-key"`
	Issuer       string `json:",default=mymall"`
}

type TelemetryConf struct {
	Enabled  bool
	Endpoint string
	Service  string `json:",default=user-service"`
}

func (c MySQLConf) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Dbname, c.Charset,
	)
}

func (c MySQLConf) ToPkg() pkgconfig.MySQLConfig {
	return pkgconfig.MySQLConfig{
		Username:        c.Username,
		Password:        c.Password,
		Host:            c.Host,
		Port:            c.Port,
		Dbname:          c.Dbname,
		Charset:         c.Charset,
		MaxOpenConns:    c.MaxOpenConns,
		MaxIdleConns:    c.MaxIdleConns,
		ConnMaxLifetime: c.ConnMaxLifetime,
		ConnMaxIdleTime: c.ConnMaxIdleTime,
	}
}

func (c TelemetryConf) ToPkg() pkgconfig.TelemetryConfig {
	return pkgconfig.TelemetryConfig{
		Enabled:  c.Enabled,
		Endpoint: c.Endpoint,
		Service:  c.Service,
	}
}

// ExpireHours converts AccessExpire seconds to hours for pkg/jwt.
func (c AuthConf) ExpireHours() int {
	if c.AccessExpire <= 0 {
		return 24
	}
	h := int(c.AccessExpire / 3600)
	if h < 1 {
		return 1
	}
	return h
}

// GRPCPort parses RpcListenOn host:port.
func (c Config) GRPCPort() int {
	parts := strings.Split(c.RpcListenOn, ":")
	if len(parts) == 0 {
		return 9090
	}
	p, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil || p == 0 {
		return 9090
	}
	return p
}

// OverlayFromEnv applies MYMALL_* overrides used by Compose/K8s.
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
	if v := strings.TrimSpace(os.Getenv("MYMALL_JWT_SECRET")); v != "" {
		c.Auth.AccessSecret = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_JWT_CONSUMER_KEY")); v != "" {
		c.Auth.ConsumerKey = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_JWT_ISSUER")); v != "" {
		c.Auth.Issuer = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_JWT_EXPIRE_HOURS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.Auth.AccessExpire = int64(n) * 3600
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TELEMETRY_ENABLED")); v != "" {
		c.Telemetry.Enabled = v == "1" || strings.EqualFold(v, "true")
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TELEMETRY_ENDPOINT")); v != "" {
		c.Telemetry.Endpoint = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_TELEMETRY_SERVICE")); v != "" {
		c.Telemetry.Service = v
	}
}
