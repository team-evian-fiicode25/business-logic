package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	if db == nil {
		panic("Attempted to retrieve database before its initialization")
	}

	return db
}

func InitDB(dsn string) error {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return err
}

func InitDBFromEnv() error {
	dsn := os.Getenv("POSTGRES_CONNECTION")

	if dsn == "" {
		return fmt.Errorf("missing environment variable: POSTGRES_CONNECTION")
	}

	return InitDB(dsn)
}
