package database

import (
	"fmt"
	"log"
	"mail-service/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	const migrationsDir = "migrations"
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	dbInstance, err := db.DB()
	if err != nil {
		return err
	}

	if err := goose.Up(dbInstance, migrationsDir); err != nil {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
