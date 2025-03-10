package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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
	var dsn string
	if envDSN := os.Getenv("DATABASE_URL"); envDSN != "" {
		dsn = envDSN
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port,
		)
	}

	var db *gorm.DB
	var err error
	const maxAttempts = 5
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		logger.Warn("Database connection failed", zap.Error(err), zap.Int("attempt", i))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logger.Error("Failed to connect to database after retries", zap.Error(err))
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil || sqlDB == nil {
		logger.Error("Underlying sql.DB is nil", zap.Error(err))
		return nil, fmt.Errorf("sql.DB is nil")
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Database ping failed", zap.Error(err))
		return nil, errors.New("database ping failed")
	}
	logger.Info("Connected to database")

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
