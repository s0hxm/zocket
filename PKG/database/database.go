package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"your-project/backend/internal/config"
)

// DB is the database connection
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Connected to database successfully")

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
		return nil, err
	}

	// Set max open connections
	sqlDB.SetMaxOpenConns(100)

	// Set max idle connections
	sqlDB.SetMaxIdleConns(10)

	DB = db
	return db, nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
