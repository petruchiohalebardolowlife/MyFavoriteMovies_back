package handlers

import (
	"myfavouritemovies/repository"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUserHandler(c *gin.Context) {
  var user structs.User
  if !utils.BindJSON(c, &user) {
      return
  }

  if err := repository.AddUser(&user); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  c.Status(http.StatusCreated)
}

func UpdateUserHandler(c *gin.Context) {
  user, errUser := utils.GetContextUser(c)
  if errUser!=nil || !utils.BindJSON(c, &user) {
      return
  }

  if err := repository.UpdateUser(user); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  c.Status(http.StatusOK)
}