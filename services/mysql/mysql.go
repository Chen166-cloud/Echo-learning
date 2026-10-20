// mysql dependence service
package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // 底层 MySQL 驱动
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// DB 全局 gorm 句柄
var DB *gorm.DB

// getEnv 读取环境变量，没设置就用默认值
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Init 连接 MySQL：自动建库 + 建立连接
func Init() {
	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASS", "1234")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	name := getEnv("DB_NAME", "echo_learning")

	// 先不带库名连上 MySQL 服务器，确保数据库存在（教学项目省去手动建库）
	serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
	raw, err := sql.Open("mysql", serverDSN)
	if err != nil {
		log.Fatalf("connect mysql server failure, err=%v", err)
	}
	defer raw.Close()
	if _, err := raw.Exec(
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", name),
	); err != nil {
		log.Fatalf("create database failure, err=%v", err)
	}

	// 连接具体数据库
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info), // 打印 SQL，方便学习
	})
	if err != nil {
		log.Fatalf("connect database failure, err=%v", err)
	}

	// 连接池参数
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("get raw connection failure, err=%v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// GetDB 获取全局 gorm 句柄
func GetDB() *gorm.DB {
	return DB
}
