package database

import (
	"fmt"
	"log"

	"github.com/collab-platform/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(host, port, user, password, dbname string) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	pgdb := &PostgresDB{DB: db}
	if err := pgdb.AutoMigrate(); err != nil {
		return nil, err
	}

	return pgdb, nil
}

func (p *PostgresDB) AutoMigrate() error {
	err := p.DB.AutoMigrate(
		&domain.User{},
		&domain.Document{},
		&domain.DocumentPermission{},
		&domain.DocumentVersion{},
		&domain.Activity{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}
	log.Println("Database migration completed successfully")
	return nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

