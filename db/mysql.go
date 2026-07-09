package db

import (
	"log"
	"time"

	"mymall/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormDB *gorm.DB

func InitGormMySQL(cfg config.MySQLConfig) {
	dsn := cfg.DSN()

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("GORM 连接 MySQL 失败：%v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取 GORM 底层 sql.DB 失败：%v", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	GormDB = db
	log.Println("GORM MySQL 连接池初始化成功")
}

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
