package mysql

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

// LoadDatabaseConfig loads the database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file on load db")
	}

	return &DatabaseConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

// SetupDatabase initializes the database connection using GORM
func ConnectDatabase() (*gorm.DB, error) {
	cfg := LoadDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
