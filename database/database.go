package database

import (
	"fmt"
	"log"
	"myfavouritemovies/structs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn:="host=localhost user=postgres password=3256 dbname=postgres port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v",err)
	}
	fmt.Println("Successfully connected to LocalDATABase on PostgreSQL!")

DB.AutoMigrate(&structs.User{},&structs.FavouriteMovie{}, &structs.FavouriteGenre{})
	return DB
}
