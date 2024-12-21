package database

import (
	"concert/internal/concert"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbSetup() (*gorm.DB, error) {
	loadEnvFromProjectRoot()
	DbName := os.Getenv("DB_Name")
	DbHost := os.Getenv("DB_Host")
	DbUser := os.Getenv("DB_User")
	DbPass := os.Getenv("DB_Password")
	DbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DbHost, DbPort, DbUser, DbName, DbPass)
	fmt.Println("Connecting to DB with:", connectionString)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connectionString,
	}),
	// &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Info),
	// }
	)

	dbping, _ := db.DB()
	if err != nil {
		fmt.Println("can't connect to the DB")
		return nil, err
	}
	if err := dbping.Ping(); err != nil {
		fmt.Println("Can't ping the DB")
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {

	if err := db.AutoMigrate(&concert.Artist{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&concert.Show{}); err != nil {
		return err
	}
	return db.AutoMigrate(&concert.Fan{})
}

func loadEnvFromProjectRoot() {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	for {
		envPath := filepath.Join(currentDir, ".env")
		err = godotenv.Load(envPath)
		if err == nil {
			fmt.Printf("Loaded .env from %s\n", envPath)
			return
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root without finding .env
			log.Println("Could not find .env file")
			break
		}
		currentDir = parentDir
	}
}
