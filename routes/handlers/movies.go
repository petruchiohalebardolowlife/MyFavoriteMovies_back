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

  user, errUser := utils.CheckContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  repository.AddFavoriteMovie(c, user.ID, input)
  c.Status(http.StatusCreated)
}

func ToggleWatchedStatusHandler(c *gin.Context) {
  var input struct {
      MovieID uint `json:"movie_id"`
  }

  user, errUser := utils.CheckContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  repository.ToggleWatchedStatus(c, user.ID, input.MovieID)
  c.Status(http.StatusOK)
}

func DeleteFavoriteMovieHandler(c *gin.Context) {
  var input struct {
      MovieID uint `json:"movie_id"`
  }

  user, errUser := utils.CheckContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  repository.DeleteFavoriteMovie(c, user.ID, input.MovieID)
  c.Status(http.StatusOK)
}