package handlers

import (
	"myfavouritemovies/repository"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddFavoriteMovieHandler(c *gin.Context) {
  var input structs.Movie

  user, errUser := utils.GetContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  if err := repository.AddFavoriteMovie(user.ID, input); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  c.Status(http.StatusCreated)
}

func ToggleWatchedStatusHandler(c *gin.Context) {
  var input struct {
      MovieID int32 `json:"movie_id"`
  }

  user, errUser := utils.GetContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  if err := repository.ToggleWatchedStatus(user.ID, input.MovieID); err != nil {
      c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
      return
  }

  c.Status(http.StatusOK)
}

func DeleteFavoriteMovieHandler(c *gin.Context) {
  var input struct {
      MovieID int32 `json:"movie_id"`
  }

  user, errUser := utils.GetContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  if err := repository.DeleteFavoriteMovie(user.ID, input.MovieID); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  c.Status(http.StatusOK)
}