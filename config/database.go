package config

import (
	"context"
	"fmt"
	"log"

	"github/com/cl0ky/e-voting-be/env"
	"github/com/cl0ky/e-voting-be/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDB(ctx context.Context) *DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		env.DBHost, env.DBPort, env.DBUser, env.DBPassword, env.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Election{},
		&models.Candidate{},
		&models.EligibleVoter{},
		&models.Vote{},
	)

	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Printf("✅ Database connected! Host: %s Port: %d DB: %s", env.DBHost, env.DBPort, env.DBName)

	return &DB{ctx: ctx, db: db}
}

func (d *DB) Instance() *gorm.DB {
	return d.db
}

func (d *DB) Close() {
	sqlDB, err := d.db.DB()
	if err != nil {
		log.Fatalf("failed to get db instance: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("⚠️ error closing DB: %v", err)
	} else {
		log.Println("🛑 Database connection closed.")
	}
}
