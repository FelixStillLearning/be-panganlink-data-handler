package config

import (
	"fmt"
	"log"

	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(cfg *Config) *gorm.DB {
	// 1. Connect without database name to create the database if it doesn't exist
	dsnNoDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort)

	dbNoName, err := gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}

	createDBCommand := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.DBName)
	if err := dbNoName.Exec(createDBCommand).Error; err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// 2. Connect with the database name to perform migrations
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-Migrate schemas (Production-ready)
	err = db.AutoMigrate(
		&model.User{},
		&model.Komoditas{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.Notification{},
		&model.MarketPrice{},
	)
	if err != nil {
		log.Printf("Migration warning: %v", err)
	}

	// Seeder for Super Admin
	seedSuperAdmin(db)
	
	// Seeder for Default Commodities
	seedKomoditas(db)

	log.Println("Database connection established")
	return db
}

func seedSuperAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := model.User{
			ID:       uuid.New().String(),
			Name:     "Super Admin",
			Email:    "admin@panganlink.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}
		db.Create(&admin)
		log.Println("Super Admin seeded automatically (admin@panganlink.com / admin123)")
	}
}

func seedKomoditas(db *gorm.DB) {
	var count int64
	db.Model(&model.Komoditas{}).Count(&count)
	if count == 0 {
		commodities := []model.Komoditas{
			{ID: "KMD-001", Nama: "Beras", Satuan: "kg", Kategori: "Beras"},
			{ID: "KMD-002", Nama: "Cabai Merah", Satuan: "kg", Kategori: "Sayuran"},
			{ID: "KMD-003", Nama: "Bawang Merah", Satuan: "kg", Kategori: "Bumbu Dapur"},
		}
		db.Create(&commodities)
		log.Println("Default commodities seeded automatically")
	}
}
