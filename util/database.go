package util

import (
	"fmt"
	"log"
	"simcomm-monolith/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const SQLSTATE_23505 = "SQLSTATE 23505"

func GetDB(cfg *config.Config) *gorm.DB {

	switch cfg.DBConfig.Driver {
	case "postgresql":
		dsn := "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=Asia/Jakarta"
		dsn = fmt.Sprintf(
			dsn,
			cfg.DBConfig.Host,
			cfg.DBConfig.Port,
			cfg.DBConfig.User,
			cfg.DBConfig.Password,
			cfg.DBConfig.DBName,
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect to database")
		}
		return db
	default:
		log.Fatalf("unknown database driver : %v", cfg.DBConfig.Driver)
	}
	return nil
}
