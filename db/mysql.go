package db

import "database/sql"
import "time"
import "log"
import "fmt"
import "github.com/spf13/viper"

var MySQL *sql.DB

type MySQLConfig struct {
	Username        int    `mapstructure:"username"`           // 数据库用户名
	Password        string `mapstructure:"password"`           // 数据库密码
	Host            string `mapstructure:"host"`               // 数据库地址
	Port            int    `mapstructure:"port"`               // 数据库端口
	Dbname          string `mapstructure:"dbname"`             // 数据库名
	Charset         string `mapstructure:"charset"`            // 字符集
	MaxOpenConns    int    `mapstructure:"max_open_conns"`     // 最大打开连接数
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`     // 最大空闲连接数
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`  // 连接最大存活时间（秒）
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"` // 连接最大空闲时间（秒）
}

// InitMySQL 整合：读取config.yaml配置 + 初始化MySQL连接池
func InitMySQL() {
	// 1. 读取根目录config.yaml配置文件
	viper.SetConfigFile("./config.yaml") // 指定配置文件路径（根目录）
	viper.SetConfigType("yaml")          // 指定配置文件格式

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取config.yaml失败：%v", err)
	}

	// 2. 将yaml中的mysql配置解析到MySQLConfig结构体
	var mysqlCfg MySQLConfig
	if err := viper.UnmarshalKey("mysql", &mysqlCfg); err != nil {
		log.Fatalf("解析mysql配置失败：%v", err)
	}

	// 3. 构建MySQL DSN连接字符串（根据配置动态拼接）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.Dbname,
		mysqlCfg.Charset,
	)

	// 4. 初始化MySQL连接池
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("初始化MySQL连接池失败：%v", err)
	}

	// 5. 配置连接池参数（从配置文件读取，无需硬编码）
	db.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	db.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(mysqlCfg.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(mysqlCfg.ConnMaxIdleTime) * time.Second)

	// 6. 验证连接是否有效
	if err := db.Ping(); err != nil {
		log.Fatalf("MySQL连接验证失败：%v", err)
	}

	// 7. 赋值给全局DB实例
	MySQL = db
	log.Println("MySQL连接池初始化成功（配置来自config.yaml）")
}

// CloseMySQL 关闭数据库连接池（项目退出时调用）
func CloseMySQL() {
	if MySQL != nil {
		if err := MySQL.Close(); err != nil {
			log.Printf("关闭MySQL连接池失败：%v", err)
		} else {
			log.Println("MySQL连接池已关闭")
		}
	}
}
