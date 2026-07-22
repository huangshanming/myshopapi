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

	Etcd EtcdConf `json:",optional"`

	MySQL     MySQLConf
	Redis     RedisConf
	RabbitMQ  RabbitMQConf
	GRPC      GRPCConf
	// AppTelemetry is our OTel init config. Named to avoid clash with rest.RestConf / ServiceConf.Telemetry.
	AppTelemetry TelemetryConf
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

type RedisConf struct {
	Host     string `json:",default=127.0.0.1"`
	Port     int    `json:",default=6379"`
	Password string
	DB       int `json:",optional"`
}

type RabbitMQConf struct {
	Host     string `json:",default=127.0.0.1"`
	Port     int    `json:",default=5672"`
	Username string `json:",default=mymall"`
	Password string `json:",default=mymall"`
	Vhost    string `json:",default=/"`
	Exchange string `json:",default=mymall.events"`
}

type GRPCConf struct {
	UserService     string `json:",default=localhost:9090"`
	CatalogService  string `json:",default=localhost:9091"`
	MerchantService string `json:",default=localhost:9092"`
}

type TelemetryConf struct {
	Enabled  bool   `json:",default=false"`
	Endpoint string `json:",optional"`
	Service  string `json:",default=order-service"`
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

func (c RedisConf) ToPkg() pkgconfig.RedisConfig {
	return pkgconfig.RedisConfig{Host: c.Host, Port: c.Port, Password: c.Password, DB: c.DB}
}

func (c RabbitMQConf) ToPkg() pkgconfig.RabbitMQConfig {
	return pkgconfig.RabbitMQConfig{
		Host: c.Host, Port: c.Port, Username: c.Username, Password: c.Password,
		Vhost: c.Vhost, Exchange: c.Exchange,
	}
}

func (c TelemetryConf) ToPkg() pkgconfig.TelemetryConfig {
	return pkgconfig.TelemetryConfig{Enabled: c.Enabled, Endpoint: c.Endpoint, Service: c.Service}
}


func (c *Config) OverlayFromEnv() {
	if v := strings.TrimSpace(os.Getenv("MYMALL_SERVER_HTTP_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.Port = n
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
	if v := strings.TrimSpace(os.Getenv("MYMALL_REDIS_HOST")); v != "" {
		c.Redis.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_REDIS_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.Redis.Port = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_REDIS_PASSWORD")); v != "" {
		c.Redis.Password = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_REDIS_DB")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.Redis.DB = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_RABBITMQ_HOST")); v != "" {
		c.RabbitMQ.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_RABBITMQ_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.RabbitMQ.Port = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_RABBITMQ_USERNAME")); v != "" {
		c.RabbitMQ.Username = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_RABBITMQ_PASSWORD")); v != "" {
		c.RabbitMQ.Password = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_RABBITMQ_VHOST")); v != "" {
		c.RabbitMQ.Vhost = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_RABBITMQ_EXCHANGE")); v != "" {
		c.RabbitMQ.Exchange = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_GRPC_USER_SERVICE")); v != "" {
		c.GRPC.UserService = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_GRPC_CATALOG_SERVICE")); v != "" {
		c.GRPC.CatalogService = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_GRPC_MERCHANT_SERVICE")); v != "" {
		c.GRPC.MerchantService = v
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
