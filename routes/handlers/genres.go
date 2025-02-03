package handlers

import (
	"myfavouritemovies/repository"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddGenresHandler(c *gin.Context) {
  var input struct {
      Genres []structs.Genre `json:"genres"`
  }

  if !utils.BindJSON(c, &input) {
      return
  }

  repository.AddGenres(c, input.Genres)
  c.Status(http.StatusCreated)
}


func AddFavoriteGenreHandler(c *gin.Context) {
    var input struct {
        GenreID uint `json:"genre_id"`
    }

    user, errUser := utils.CheckContextUser(c)
    if !errUser || !utils.BindJSON(c, &input) {
        return
    }

    repository.AddFavoriteGenre(c, user.ID, input.GenreID)
    c.Status(http.StatusCreated)
}

func DeleteFavoriteGenreHandler(c *gin.Context) {
  var input struct {
      GenreID uint `json:"genre_id"`
  }

  user, errUser := utils.CheckContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  repository.DeleteFavoriteGenre(c, user.ID, input.GenreID)
  c.Status(http.StatusOK)
}
