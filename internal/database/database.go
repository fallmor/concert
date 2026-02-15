package database

import (
	"concert/internal/models"
	"concert/internal/utils"
	"fmt"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbSetup() (*gorm.DB, error) {
	projectRoot := utils.GetProjectRoot("go.mod")
	envPath := filepath.Join(projectRoot, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}
	DbName := utils.GetEnvOrDefault("DB_Name", "")
	DbHost := utils.GetEnvOrDefault("DB_Host", "")
	DbUser := utils.GetEnvOrDefault("DB_User", "")
	DbPass := utils.GetEnvOrDefault("DB_Password", "")
	DbPort := utils.GetEnvOrDefault("DB_PORT", "")
	DbSSLMode := utils.GetEnvOrDefault("DB_SSLMode", "")
	if DbSSLMode == "" {
		if DbHost == "localhost" || DbHost == "127.0.0.1" {
			DbSSLMode = "disable"
		} else {
			DbSSLMode = "require"
		}
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", DbHost, DbPort, DbUser, DbName, DbPass, DbSSLMode)
	log.Printf("Connecting to database at %s:%s", DbHost, DbPort)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connectionString,
	}))
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	dbPing, err := db.DB()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return nil, err
	}

	if err := dbPing.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Artist{},
		&models.Show{},
		&models.Booking{}, 
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migration completed")
	return nil
}
