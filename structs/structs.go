package structs

import (
	"gorm.io/gorm"
)

type User struct {
  gorm.Model
	ID             int32           `json:"id"`
	NickName       string           `json:"nickName"`
	UserName       string           `json:"userName"`
  Password string `json:"password"`
	FavoriteMovies []*FavoriteMovie `json:"favoriteMovies"`
	FavoriteGenres []*FavoriteGenre `json:"favoriteGenres"`
}


type Genre struct {
  gorm.Model
  ID   int32   `gorm:"primaryKey" json:"id"`
  Name string `json:"name"`
}

type FavoriteMovie struct {
  gorm.Model
  UserID      int32    `json:"user_id"`
  MovieID     int32    `json:"movie_id"`
  Title       string  `json:"title"`
  PosterPath  string  `json:"poster_path"`
  VoteAverage float64 `json:"vote_average"`
  Watched     bool    `json:"watched"`
  Genres      []*Genre `gorm:"many2many:favorite_movie_genres;" json:"genres"`
}

type FavoriteGenre struct {
  gorm.Model
  UserID  int32
  GenreID int32
}

type SignIn struct {
  Username string `json:"username"`
  Password string `json:"password"`
}
