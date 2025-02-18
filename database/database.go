package database

import (
	"fmt"
	"log"
	config "myfavouritemovies/configs"
	"myfavouritemovies/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
  config.LoadConfig()
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
    config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_PORT, config.DB_SSLMODE)
  
  db, err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
  if err != nil {
    log.Fatal("Failed to connect to database ",err)
  }
  DB=db
  fmt.Println("Successfully connected to LocalDATABase on PostgreSQL!")
  DB.AutoMigrate(&models.User{},&models.FavoriteMovie{}, &models.FavoriteGenre{}, &models.Genre{}, &models.Session{}, &models.BlackListToken{})
  
  return DB
}
