package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"wallet/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(config *Config) (*gorm.DB, error) {
	// 构建DSN (数据源名称)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.DBName,
		config.MySQL.Charset,
	)
	// 配置GORM日志级别
	logLevel := logger.Info
	if config.Log.Level == "debug" {
		logLevel = logger.Info
	} else if config.Log.Level == "warn" {
		logLevel = logger.Warn
	} else if config.Log.Level == "error" {
		logLevel = logger.Error
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到标准输出
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  logLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录不存在错误
			Colorful:                  true,        // 彩色输出
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database object: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间

	log.Println("Database connection established successfully")

	// 设置全局DB变量
	DB = db

	// 自动迁移表结构
	if err := db.AutoMigrate(&models.Users{}, &models.Wallets{}, &models.Transaction{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}
