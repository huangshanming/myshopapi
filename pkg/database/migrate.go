package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

// ShouldAutoMigrate：debug 模式默认开启；生产 release 关闭。
// 可用 MYMALL_MYSQL_AUTO_MIGRATE=true|false 强制覆盖。
func ShouldAutoMigrate(serverMode string) bool {
	if v := strings.TrimSpace(os.Getenv("MYMALL_MYSQL_AUTO_MIGRATE")); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	return strings.EqualFold(serverMode, "debug")
}

// AutoMigrateIfDebug 类似 Beego RunSyncdb：按 model 建表/加列（不会删列、不会改复杂索引约束）。
// 仅本地 debug（或显式 env）启用；生产请用 scripts/*.sql。
func AutoMigrateIfDebug(serverMode string, db *gorm.DB, models ...interface{}) error {
	if !ShouldAutoMigrate(serverMode) {
		return nil
	}
	if len(models) == 0 {
		return nil
	}
	log.Printf("GORM AutoMigrate enabled (mode=%s)", serverMode)
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	log.Println("GORM AutoMigrate 完成")
	return nil
}
