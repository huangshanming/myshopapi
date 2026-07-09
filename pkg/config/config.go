package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type ServerConfig struct {
	HTTPPort int    `mapstructure:"http_port"`
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

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetEnvPrefix("MYMALL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

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
}
