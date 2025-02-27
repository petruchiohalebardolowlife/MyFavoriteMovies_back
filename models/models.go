package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
  gorm.Model
  ID             uint             `json:"id" gorm:"primaryKey"`
  NickName       string           `json:"nickName" `
  UserName       string           `json:"userName" gorm:"unique"`
  PasswordHash   string           `json:"passwordhash" gorm:"not null"`
  FavoriteMovies []*FavoriteMovie `json:"favoriteMovies" gorm:"constraint:OnDelete:CASCADE;"`
  FavoriteGenres []*FavoriteGenre `json:"favoriteGenres" gorm:"constraint:OnDelete:CASCADE;"`
}

type Genre struct {
  gorm.Model
  ID   uint   `gorm:"primaryKey" json:"id"`
  Name string `json:"name"`
}

type FavoriteMovie struct {
  gorm.Model
  ID          uint     `json:"id" gorm:"primaryKey"`
  UserID      uint     `json:"user_id"`
  MovieID     uint     `json:"movie_id"`
  Title       string   `json:"title"`
  PosterPath  string   `json:"poster_path"`
  VoteAverage float64  `json:"vote_average"`
  Watched     bool     `json:"watched"`
  Genres      []*Genre `gorm:"many2many:favorite_movie_genres;" json:"genres"`
}

type FavoriteGenre struct {
  gorm.Model
  UserID  uint
  GenreID uint
}

type SignIn struct {
  Username uint `json:"username"`
  Password uint `json:"password"`
}

type MovieDetails struct {
  Title       string   `json:"title"`
  Rating      float64  `json:"vote_average"`
  ReleaseDate string   `json:"release_date"`
  PosterPath  string   `json:"poster_path"`
  Genres      []*Genre `json:"genres"`
  Overview    string   `json:"overview"`
}

type Movie struct {
  ID          uint    `json:"id"`
  Title       string  `json:"title"`
  PosterPath  string  `json:"poster_path"`
  VoteAverage float64 `json:"vote_average"`
  GenreIDs    []uint  `json:"genre_ids"`
  ReleaseDate string  `json:"release_date"`
}

type Session struct {
  gorm.Model
  ID        string `gorm:"primary_key"`
  UserID    uint
  ExpiresAt time.Time
}

type BlackListToken struct {
  gorm.Model
  ID        string `gorm:"primary_key"`
  UserID    uint
  ExpiresAt time.Time
}

type TokenClaims struct {
  jwt.RegisteredClaims
  UserID uint `json:"user_id"`
}

type Token struct {
  Value  string
  Claims *TokenClaims
}

type Tokens struct {
  Access  *Token
  Refresh *Token
}
