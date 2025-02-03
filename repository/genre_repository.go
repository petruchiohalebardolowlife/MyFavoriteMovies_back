package repository

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func GetAllGenres() ([]structs.Genre, error) {
	var genres []structs.Genre
	if err := database.DB.Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}


func AddGenres(c *gin.Context, genres []structs.Genre) {
  existingGenres, err := GetAllGenres()
  if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  var newGenres []structs.Genre
  for _, genre := range genres {
      if !slices.ContainsFunc(existingGenres, func(gen structs.Genre) bool {
          return gen.ID == genre.ID
      }) {
          newGenres = append(newGenres, genre)
      }
  }

  if len(newGenres) > 0 {
      if err := database.DB.Create(&newGenres).Error; err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
          return
      }
  }
}

func AddFavoriteGenre(c *gin.Context, userID, genreID uint) {
  if err := database.DB.Where("user_id = ? AND genre_id = ?", userID, genreID).First(&structs.FavoriteGenre{}).Error; err == nil {
      c.JSON(http.StatusConflict, gin.H{"error": "genre already in favorites"})
      return
  }

  newFavorite := structs.FavoriteGenre{
      UserID:  userID,
      GenreID: genreID,
  }

  if err := database.DB.Create(&newFavorite).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  c.Status(http.StatusCreated)
}

func DeleteFavoriteGenre(c *gin.Context, userID, genreID uint) {
  var favGenre structs.FavoriteGenre
  if err := database.DB.Where("user_id = ? AND genre_id = ?", userID, genreID).First(&favGenre).Error; err != nil {
      c.JSON(http.StatusConflict, gin.H{"error": "genre not in favorites"})
      return
  }

  if err := database.DB.Delete(&favGenre).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  c.Status(http.StatusOK)
}