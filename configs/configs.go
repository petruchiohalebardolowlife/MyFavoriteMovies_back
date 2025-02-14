package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
  DB_HOST     string
  DB_USER     string
  DB_PASSWORD string
  DB_NAME     string
  DB_PORT     string
  DB_SSLMODE  string
  SRVR_PORT   string
  API_KEY     string
  TMDB_API_BASE_URL string
  APP_ENV     string
  SECRET      string
)

func LoadConfig() {
  err := godotenv.Load()
  if err != nil {
    panic("Error loading .env file")
  }

  DB_HOST = os.Getenv("DB_HOST")
  DB_USER = os.Getenv("DB_USER")
  DB_PASSWORD = os.Getenv("DB_PASSWORD")
  DB_NAME = os.Getenv("DB_NAME")
  DB_PORT = os.Getenv("DB_PORT")
  DB_SSLMODE = os.Getenv("DB_SSLMODE")
  SRVR_PORT = os.Getenv("SRVR_PORT")
  API_KEY = os.Getenv("API_KEY")
  TMDB_API_BASE_URL = os.Getenv("TMDB_API_BASE_URL")
  APP_ENV=os.Getenv("APP_ENV")
  SECRET=os.Getenv("SECRET")
}