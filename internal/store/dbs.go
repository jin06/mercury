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

var Default *gorm.DB

func Init() error {
	db, err := NewClient(&config.Def.Database)
	if err != nil {
		return err
	}
	Default = db
	return nil
}

func NewClient(cfg *config.Database) (db *gorm.DB, err error) {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)

	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Type)
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
