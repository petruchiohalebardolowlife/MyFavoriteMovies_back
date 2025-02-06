package handlers

import (
	"myfavouritemovies/repository"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddFavoriteGenreHandler(c *gin.Context) {
    var input struct {
        GenreID uint `json:"genre_id"`
    }

    user, errUser := utils.GetContextUser(c)
    if !errUser || !utils.BindJSON(c, &input) {
        return
    }

    if err := repository.AddFavoriteGenre(uint(user.ID), input.GenreID); err != nil {
        if err.Error() == "genre already in favorites" {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusCreated)
}

func DeleteFavoriteGenreHandler(c *gin.Context) {
  var input struct {
      GenreID uint `json:"genre_id"`
  }

  user, errUser := utils.GetContextUser(c)
  if !errUser || !utils.BindJSON(c, &input) {
      return
  }

  if err := repository.DeleteFavoriteGenre(uint(user.ID), input.GenreID); err != nil {
      if err.Error() == "genre not in favorites" {
          c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
          return
      }
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
  c.Status(http.StatusOK)
}
