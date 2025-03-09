package database

import (
	"fmt"
	"log"
	"mail-service/config"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func InitDB(cfg *config.Config, logger *zap.Logger) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	log.Println("Connected to database")

	// Автоматически применяем миграции
	if err := RunMigrations(db); err != nil {
		logger.Error("Failed to apply migrations", zap.Error(err))
		return nil, err
	}

	return &Database{DB: db}, nil
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

func (d *Database) Close() {
	dbInstance, err := d.DB.DB()
	if err != nil {
		log.Println("Failed to get DB instance:", err)
		return
	}
	dbInstance.Close()
}
