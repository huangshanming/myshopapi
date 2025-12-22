package db

import (
	"log"
	"time"

	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局 GORM DB 实例（替代原来的 *sql.DB）
var GormDB *gorm.DB

// MySQLConfig 仍复用原有结构体，与 config.yaml 对应
type MySQLConfig struct {
	Username        string `mapstructure:"username"` // 注意：用户名应为 string 类型（之前代码中是 int，需修正！）
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

// InitGormMySQL 初始化 GORM MySQL 连接
func InitGormMySQL() {
	// 1. 读取 config.yaml 配置
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取 config.yaml 失败：%v", err)
	}

	// 2. 解析配置到结构体
	var mysqlCfg MySQLConfig
	if err := viper.UnmarshalKey("mysql", &mysqlCfg); err != nil {
		log.Fatalf("解析 mysql 配置失败：%v", err)
	}

	// 3. 构建 GORM MySQL DSN（格式与原生一致）
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn,
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.Dbname,
		mysqlCfg.Charset,
	)

	// 4. 初始化 GORM 连接
	gormConfig := &gorm.Config{
		// 日志配置：开发环境显示 SQL 语句，生产环境可关闭
		Logger: logger.Default.LogMode(logger.Info),
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("GORM 连接 MySQL 失败：%v", err)
	}

	// 5. 配置连接池（与原生 sql.DB 一致，GORM 底层复用了 sql.DB）
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取 GORM 底层 sql.DB 失败：%v", err)
	}
	// 设置连接池参数
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlCfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(mysqlCfg.ConnMaxIdleTime) * time.Second)

	// 6. 赋值全局 GORM 实例
	GormDB = db
	log.Println("GORM MySQL 连接池初始化成功（配置来自 config.yaml）")
}

// CloseGormMySQL 关闭 GORM 连接池
func CloseGormMySQL() {
	if GormDB != nil {
		sqlDB, err := GormDB.DB()
		if err != nil {
			log.Printf("获取 GORM 底层 sql.DB 失败：%v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("关闭 GORM MySQL 连接池失败：%v", err)
		} else {
			log.Println("GORM MySQL 连接池已关闭")
		}
	}
}
