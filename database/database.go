package database

import (
	"fmt"
	"log"
	"myfavouritemovies/structs"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	err:=godotenv.Load()
	if err != nil {
		log.Println("Not found .env file")
	}
	host:=os.Getenv("DB_HOST")
	user:=os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
	
	db, err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database ",err)
	}
	DB=db
	fmt.Println("Successfully connected to LocalDATABase on PostgreSQL!")
DB.AutoMigrate(&structs.User{},&structs.FavoriteMovie{}, &structs.FavoriteGenre{})
	return DB
}
