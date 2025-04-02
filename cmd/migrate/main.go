package main

import (
	"fmt"
	"log"
	"os"

	"github.com/team-evian-fiicode25/business-logic/data"
	"github.com/team-evian-fiicode25/business-logic/database"

	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("POSTGRES_CONNECTION")

	if dsn == "" {
		log.Fatalln("Missing env variable: POSTGRES_CONNECTION")
	}

	err := database.InitDB(dsn)

	if err != nil {
		log.Fatal(err)
	}

	var db *gorm.DB = database.GetDB()

	collections := [...]any{&data.User{}}

	for _, collection := range collections {
		err := db.AutoMigrate(collection)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Completed successfully")
}
