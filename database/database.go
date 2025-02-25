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

func InitDB(host, port, user, password, dbname, timezone string) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbname, port, timezone)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return err
}

func InitDBFromEnv() error {
	var envNames = [6]string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DATABASE",
		"POSTGRES_TIMEZONE",
	}

	var values [6]string

	for i, name := range envNames {
		values[i] = os.Getenv(name)

		if values[i] == "" {
			return fmt.Errorf("missing required environment variable: %s", name)
		}
	}

	err := InitDB(values[0], values[1], values[2], values[3], values[4], values[5])
	if err != nil {
		return err
	}

	return nil
}
