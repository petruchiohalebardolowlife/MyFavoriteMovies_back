package structs

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uint `gorm:"primaryKey" json:"id"`
	Login string `gorm:"unique"`
	Username  string `gorm:"unique"`
	Pass  string `gorm:"not null"`
}

type Genre struct {
	gorm.Model
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `json:"name"`
}

type FavoriteMovie struct {
    gorm.Model
    UserID    uint  `json:"user_id"`
    MovieID   uint  `json:"movie_id"`
	Title      string `json:"title"`
	PosterPath  string `json:"poster_path"`
	VoteAverage float64 `json:"vote_average"`
    Watched   bool  `json:"watched"`
    User      User  `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID" json:"user"`
}

type FavoriteGenre struct {
    gorm.Model
    UserID  uint `json:"user_id"`
    GenreID uint `json:"genre_id"`
    User    User  `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID" json:"user"`
    Genre   Genre `gorm:"constraint:OnDelete:CASCADE;foreignKey:GenreID" json:"genre"`
}