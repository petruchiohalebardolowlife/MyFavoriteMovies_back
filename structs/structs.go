package structs

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uint `gorm:"primaryKey" json:"id"`
	Login string `gorm:"unique"`
	Username  string `gorm:"unique"`
	Pass  string `gorm:"not null"`
}

type Movie struct {
	gorm.Model
	ID          uint `gorm:"primaryKey" json:"id"`
	Tittle      string `json:"title"`
	PosterPath  string `json:"poster_path"`
	Genre_IDs pq.Int64Array `json:"genre_ids" gorm:"type:integer[]"` 
	ReleaseDate string `json:"release_date"`
	VoteAverage float64 `json:"vote_average"`
}

type Genre struct {
	gorm.Model
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `json:"name"`
}

type FavouriteMovie struct {
    gorm.Model
    UserID    uint  `json:"user_id"`
    MovieID   uint  `json:"movie_id"`
    Watched   bool  `json:"watched"`
    User      User  `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID" json:"user"`
    Movie     Movie `gorm:"constraint:OnDelete:CASCADE;foreignKey:MovieID" json:"movie"`
}

type FavouriteGenre struct {
    gorm.Model
    UserID  uint `json:"user_id"`
    GenreID uint `json:"genre_id"`
    User    User  `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID" json:"user"`
    Genre   Genre `gorm:"constraint:OnDelete:CASCADE;foreignKey:GenreID" json:"genre"`
}