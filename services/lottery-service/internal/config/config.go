package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	MySQL           MySQLConf
	Redis           RedisConf
	UserHTTP        string `json:",default=http://127.0.0.1:8881"`
	UserHTTPTimeout int    `json:",default=8"`
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
	Password string `json:",optional"`
	DB       int    `json:",default=0"`
}

func (c MySQLConf) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Dbname, c.Charset,
	)
}

func (c RedisConf) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
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
	if v := strings.TrimSpace(os.Getenv("MYMALL_REDIS_HOST")); v != "" {
		c.Redis.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_REDIS_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.Redis.Port = n
		}
	}
	if v, ok := os.LookupEnv("MYMALL_REDIS_PASSWORD"); ok {
		c.Redis.Password = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_REDIS_DB")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			c.Redis.DB = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_USER_HTTP")); v != "" {
		c.UserHTTP = v
	}
	if v := strings.TrimSpace(os.Getenv("MYMALL_USER_HTTP_TIMEOUT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.UserHTTPTimeout = n
		}
	}
}
