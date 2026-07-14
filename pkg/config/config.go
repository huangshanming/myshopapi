package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
}

type ServerConfig struct {
	HTTPPort int    `mapstructure:"http_port"`
	GRPCPort int    `mapstructure:"grpc_port"`
	Mode     string `mapstructure:"mode"`
}

type MySQLConfig struct {
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Dbname          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ConsumerKey string `mapstructure:"consumer_key"`
	ExpireHours int    `mapstructure:"expire_hours"`
	Issuer      string `mapstructure:"issuer"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Vhost    string `mapstructure:"vhost"`
	Exchange string `mapstructure:"exchange"`
}

type GRPCConfig struct {
	UserService    string `mapstructure:"user_service"`
	CatalogService string `mapstructure:"catalog_service"`
}

type TelemetryConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
	Service  string `mapstructure:"service"`
}

func (c MySQLConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Dbname,
		c.Charset,
	)
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.Username, c.Password, c.Host, c.Port, strings.TrimPrefix(c.Vhost, "/"))
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetEnvPrefix("MYMALL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	// Unmarshal 不会自动吃嵌套 AutomaticEnv，需显式 BindEnv（K8s Secret / Compose 环境变量）
	for _, key := range []string{
		"server.http_port", "server.grpc_port", "server.mode",
		"mysql.host", "mysql.port", "mysql.username", "mysql.password", "mysql.dbname", "mysql.charset",
		"jwt.secret", "jwt.consumer_key", "jwt.expire_hours", "jwt.issuer",
		"redis.host", "redis.port", "redis.password", "redis.db",
		"rabbitmq.host", "rabbitmq.port", "rabbitmq.username", "rabbitmq.password", "rabbitmq.vhost", "rabbitmq.exchange",
		"grpc.user_service", "grpc.catalog_service",
		"telemetry.enabled", "telemetry.endpoint", "telemetry.service",
	} {
		_ = v.BindEnv(key)
	}

	if path != "" {
		if _, err := os.Stat(path); err == nil {
			v.SetConfigFile(path)
			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("read config: %w", err)
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.Server.HTTPPort == 0 {
		c.Server.HTTPPort = 8087
	}
	if c.Server.GRPCPort == 0 {
		c.Server.GRPCPort = 9090
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}
	if c.MySQL.Charset == "" {
		c.MySQL.Charset = "utf8mb4"
	}
	if c.MySQL.MaxOpenConns == 0 {
		c.MySQL.MaxOpenConns = 100
	}
	if c.MySQL.MaxIdleConns == 0 {
		c.MySQL.MaxIdleConns = 16
	}
	if c.MySQL.ConnMaxLifetime == 0 {
		c.MySQL.ConnMaxLifetime = 3600
	}
	if c.MySQL.ConnMaxIdleTime == 0 {
		c.MySQL.ConnMaxIdleTime = 3600
	}
	if c.JWT.ExpireHours == 0 {
		c.JWT.ExpireHours = 24
	}
	if c.JWT.Issuer == "" {
		c.JWT.Issuer = "mymall"
	}
	if c.JWT.ConsumerKey == "" {
		c.JWT.ConsumerKey = "mymall-user-key"
	}
	if c.Redis.Port == 0 {
		c.Redis.Port = 6379
	}
	if c.RabbitMQ.Port == 0 {
		c.RabbitMQ.Port = 5672
	}
	if c.RabbitMQ.Vhost == "" {
		c.RabbitMQ.Vhost = "/"
	}
	if c.RabbitMQ.Exchange == "" {
		c.RabbitMQ.Exchange = "mymall.events"
	}
	if c.GRPC.UserService == "" {
		c.GRPC.UserService = "localhost:9090"
	}
	if c.GRPC.CatalogService == "" {
		c.GRPC.CatalogService = "localhost:9091"
	}
}
