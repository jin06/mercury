package store

import (
	"fmt"
	"log"
	"time"

	"github.com/jin06/mercury/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewClient(cfg *config.DBConfig) (db *gorm.DB, err error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		dialector = mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
		dialector = postgres.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// GORM 配置
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	log.Println("Database connection established successfully.")
	return db, nil
}
