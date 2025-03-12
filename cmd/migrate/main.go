package main

import (
	"fmt"
	"log"

	"github.com/team-evian-fiicode25/business-logic/data"
	"github.com/team-evian-fiicode25/business-logic/database"

	"gorm.io/gorm"
)

func main() {
    err := database.InitDBFromEnv();

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
