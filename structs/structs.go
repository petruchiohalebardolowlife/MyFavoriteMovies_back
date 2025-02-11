package structs

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uint             `json:"id" gorm:"primaryKey"`
	NickName       string           `json:"nickName" `
	UserName       string           `json:"userName" gorm:"unique"`
	Password       string           `json:"password" gorm:"not null"`
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
