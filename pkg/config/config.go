// Package config is Viper-based bootstrap for inventory-sync-service only.
// Business services (user/catalog/order/merchant) use go-zero conf.MustLoad + rest.RestConf.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig    `mapstructure:"server"`
	MySQL        MySQLConfig     `mapstructure:"mysql"`
	JWT          JWTConfig       `mapstructure:"jwt"`
	Redis        RedisConfig     `mapstructure:"redis"`
	RabbitMQ     RabbitMQConfig  `mapstructure:"rabbitmq"`
	GRPC         GRPCConfig      `mapstructure:"grpc"`
	MerchantHTTP string          `mapstructure:"merchant_http"`
	UserHTTP     string          `mapstructure:"user_http"`
	Telemetry    TelemetryConfig `mapstructure:"telemetry"`
	Canal        CanalConfig     `mapstructure:"canal"`
}

// CanalConfig is used by inventory-sync-service (canal-go client).
type CanalConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Destination string `mapstructure:"destination"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Filter      string `mapstructure:"filter"`
	SoTimeout   int    `mapstructure:"so_timeout"`
	IdleTimeout int    `mapstructure:"idle_timeout"`
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
		"merchant_http", "user_http",
		"telemetry.enabled", "telemetry.endpoint", "telemetry.service",
		"canal.host", "canal.port", "canal.destination", "canal.username", "canal.password", "canal.filter",
		"canal.so_timeout", "canal.idle_timeout",
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
	// viper.Unmarshal 对嵌套 key 不会应用 BindEnv/AutomaticEnv，需用 Get* 再盖一层
	overlayFromViper(&cfg, v)

	cfg.applyDefaults()
	return &cfg, nil
}

func overlayFromViper(cfg *Config, v *viper.Viper) {
	// Prefer explicit env（K8s Secret / Compose）；viper.Unmarshal 不会应用嵌套 AutomaticEnv
	envOr := func(envKey, viperKey string) string {
		if s := strings.TrimSpace(os.Getenv(envKey)); s != "" {
			return s
		}
		return v.GetString(viperKey)
	}
	envIntOr := func(envKey, viperKey string, cur int) int {
		if s := strings.TrimSpace(os.Getenv(envKey)); s != "" {
			var n int
			if _, err := fmt.Sscanf(s, "%d", &n); err == nil && n != 0 {
				return n
			}
		}
		if n := v.GetInt(viperKey); n != 0 {
			return n
		}
		return cur
	}

	if s := envOr("MYMALL_SERVER_MODE", "server.mode"); s != "" {
		cfg.Server.Mode = s
	}
	cfg.Server.HTTPPort = envIntOr("MYMALL_SERVER_HTTP_PORT", "server.http_port", cfg.Server.HTTPPort)
	cfg.Server.GRPCPort = envIntOr("MYMALL_SERVER_GRPC_PORT", "server.grpc_port", cfg.Server.GRPCPort)

	if s := envOr("MYMALL_MYSQL_HOST", "mysql.host"); s != "" {
		cfg.MySQL.Host = s
	}
	if s := envOr("MYMALL_MYSQL_USERNAME", "mysql.username"); s != "" {
		cfg.MySQL.Username = s
	}
	if s := envOr("MYMALL_MYSQL_PASSWORD", "mysql.password"); s != "" {
		cfg.MySQL.Password = s
	}
	if s := envOr("MYMALL_MYSQL_DBNAME", "mysql.dbname"); s != "" {
		cfg.MySQL.Dbname = s
	}
	if s := envOr("MYMALL_MYSQL_CHARSET", "mysql.charset"); s != "" {
		cfg.MySQL.Charset = s
	}
	cfg.MySQL.Port = envIntOr("MYMALL_MYSQL_PORT", "mysql.port", cfg.MySQL.Port)

	if s := envOr("MYMALL_JWT_SECRET", "jwt.secret"); s != "" {
		cfg.JWT.Secret = s
	}
	if s := envOr("MYMALL_JWT_CONSUMER_KEY", "jwt.consumer_key"); s != "" {
		cfg.JWT.ConsumerKey = s
	}
	if s := envOr("MYMALL_JWT_ISSUER", "jwt.issuer"); s != "" {
		cfg.JWT.Issuer = s
	}
	cfg.JWT.ExpireHours = envIntOr("MYMALL_JWT_EXPIRE_HOURS", "jwt.expire_hours", cfg.JWT.ExpireHours)

	if s := envOr("MYMALL_REDIS_HOST", "redis.host"); s != "" {
		cfg.Redis.Host = s
	}
	if s := envOr("MYMALL_REDIS_PASSWORD", "redis.password"); s != "" {
		cfg.Redis.Password = s
	}
	cfg.Redis.Port = envIntOr("MYMALL_REDIS_PORT", "redis.port", cfg.Redis.Port)
	if s := strings.TrimSpace(os.Getenv("MYMALL_REDIS_DB")); s != "" {
		var n int
		if _, err := fmt.Sscanf(s, "%d", &n); err == nil {
			cfg.Redis.DB = n
		}
	} else if v.IsSet("redis.db") {
		cfg.Redis.DB = v.GetInt("redis.db")
	}

	if s := envOr("MYMALL_RABBITMQ_HOST", "rabbitmq.host"); s != "" {
		cfg.RabbitMQ.Host = s
	}
	if s := envOr("MYMALL_RABBITMQ_USERNAME", "rabbitmq.username"); s != "" {
		cfg.RabbitMQ.Username = s
	}
	if s := envOr("MYMALL_RABBITMQ_PASSWORD", "rabbitmq.password"); s != "" {
		cfg.RabbitMQ.Password = s
	}
	if s := envOr("MYMALL_RABBITMQ_VHOST", "rabbitmq.vhost"); s != "" {
		cfg.RabbitMQ.Vhost = s
	}
	if s := envOr("MYMALL_RABBITMQ_EXCHANGE", "rabbitmq.exchange"); s != "" {
		cfg.RabbitMQ.Exchange = s
	}
	cfg.RabbitMQ.Port = envIntOr("MYMALL_RABBITMQ_PORT", "rabbitmq.port", cfg.RabbitMQ.Port)

	if s := envOr("MYMALL_GRPC_USER_SERVICE", "grpc.user_service"); s != "" {
		cfg.GRPC.UserService = s
	}
	if s := envOr("MYMALL_GRPC_CATALOG_SERVICE", "grpc.catalog_service"); s != "" {
		cfg.GRPC.CatalogService = s
	}
	if s := envOr("MYMALL_MERCHANT_HTTP", "merchant_http"); s != "" {
		cfg.MerchantHTTP = s
	}
	if s := envOr("MYMALL_USER_HTTP", "user_http"); s != "" {
		cfg.UserHTTP = s
	}

	if s := strings.TrimSpace(os.Getenv("MYMALL_TELEMETRY_ENABLED")); s != "" {
		cfg.Telemetry.Enabled = s == "1" || strings.EqualFold(s, "true")
	} else if v.IsSet("telemetry.enabled") {
		cfg.Telemetry.Enabled = v.GetBool("telemetry.enabled")
	}
	if s := envOr("MYMALL_TELEMETRY_ENDPOINT", "telemetry.endpoint"); s != "" {
		cfg.Telemetry.Endpoint = s
	}
	if s := envOr("MYMALL_TELEMETRY_SERVICE", "telemetry.service"); s != "" {
		cfg.Telemetry.Service = s
	}

	if s := envOr("MYMALL_CANAL_HOST", "canal.host"); s != "" {
		cfg.Canal.Host = s
	}
	cfg.Canal.Port = envIntOr("MYMALL_CANAL_PORT", "canal.port", cfg.Canal.Port)
	if s := envOr("MYMALL_CANAL_DESTINATION", "canal.destination"); s != "" {
		cfg.Canal.Destination = s
	}
	if s := envOr("MYMALL_CANAL_USERNAME", "canal.username"); s != "" {
		cfg.Canal.Username = s
	}
	if s := envOr("MYMALL_CANAL_PASSWORD", "canal.password"); s != "" {
		cfg.Canal.Password = s
	}
	if s := envOr("MYMALL_CANAL_FILTER", "canal.filter"); s != "" {
		cfg.Canal.Filter = s
	}
	cfg.Canal.SoTimeout = envIntOr("MYMALL_CANAL_SO_TIMEOUT", "canal.so_timeout", cfg.Canal.SoTimeout)
	cfg.Canal.IdleTimeout = envIntOr("MYMALL_CANAL_IDLE_TIMEOUT", "canal.idle_timeout", cfg.Canal.IdleTimeout)
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
	if c.MerchantHTTP == "" {
		c.MerchantHTTP = "http://127.0.0.1:8884"
	}
	if c.UserHTTP == "" {
		c.UserHTTP = "http://127.0.0.1:8881"
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
	if c.Canal.Host == "" {
		c.Canal.Host = "127.0.0.1"
	}
	if c.Canal.Port == 0 {
		c.Canal.Port = 11111
	}
	if c.Canal.Destination == "" {
		c.Canal.Destination = "mymall"
	}
	if c.Canal.Filter == "" {
		c.Canal.Filter = `mymall\.(product_skus|lottery_prizes)`
	}
	if c.Canal.SoTimeout == 0 {
		c.Canal.SoTimeout = 60000
	}
	if c.Canal.IdleTimeout == 0 {
		c.Canal.IdleTimeout = 60 * 60 * 1000
	}
}
