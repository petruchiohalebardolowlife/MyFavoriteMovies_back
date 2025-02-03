package repository

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddFavoriteMovie(c *gin.Context, userID uint, input structs.Movie) {
  if _, err := utils.FindFavoriteMovie(userID, input.MovieID); err == nil {
      c.JSON(http.StatusConflict, gin.H{"error": "Movie already in favorites"})
      return
  }

  newFavorite := structs.FavoriteMovie{
      UserID:      userID,
      MovieID:     input.MovieID,
      Title:       input.Title,
      PosterPath:  input.PosterPath,
      VoteAverage: input.VoteAverage,
      Genres:      input.Genres,
      Watched:     false,
  }

  if err := database.DB.Create(&newFavorite).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  if err := database.DB.Model(&newFavorite).Association("Genres").Append(input.Genres); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate genres"})
      return
  }
}

func ToggleWatchedStatus(c *gin.Context, userID, movieID uint) {
  favMovie, err := utils.FindFavoriteMovie(userID, movieID)
  if err != nil {
      c.Status(http.StatusNotFound)
      return
  }

  favMovie.Watched = !favMovie.Watched

  if err := database.DB.Save(&favMovie).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
}

func DeleteFavoriteMovie(c *gin.Context, userID, movieID uint) {
  existingMovie, err := utils.FindFavoriteMovie(userID, movieID)
  if err != nil {
      c.Status(http.StatusNotFound)
      return
  }

  if err := database.DB.Delete(&existingMovie).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
}