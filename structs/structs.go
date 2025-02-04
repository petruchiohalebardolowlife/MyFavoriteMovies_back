package structs

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uint `gorm:"primaryKey" json:"id"`
	NickName string `gorm:"unique" json:"nick_name"` 
	Username  string `gorm:"unique" json:"username"`
	Password  string `gorm:"not null" json:"password"`
    FavoriteMovies  []FavoriteMovie  `gorm:"constraint:OnDelete:CASCADE;" json:"favorite_movies"`
    FavoriteGenres []FavoriteGenre `gorm:"constraint:OnDelete:CASCADE;" json:"favorite_genres"`
}

type Genre struct {
	gorm.Model
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `json:"name"`
}

type FavoriteMovie struct {
    gorm.Model
    UserID      uint `json:"user_id"`
    MovieID     uint  `json:"movie_id"`
    Title       string `json:"title"`
    PosterPath  string  `json:"poster_path"`
    VoteAverage float64 `json:"vote_average"`
    Watched     bool    `json:"watched"`
    Genres      []Genre `gorm:"many2many:favorite_movie_genres;" json:"genres"`
}

type FavoriteGenre struct {
    gorm.Model
    UserID  uint
    GenreID uint
}

type Movie struct {
    MovieID     uint    `json:"id"`
    Title       string  `json:"title"`
    PosterPath  string  `json:"poster_path"`
    VoteAverage float64 `json:"vote_average"`
    GenreIDs    []uint `json:"genre_ids"`
    Genres []Genre `json:"genres"`
}

type SignIn struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type MovieFilter struct {
    GenreIDs []int
    Rating   string 
    Year     int  
    Page     uint 
}