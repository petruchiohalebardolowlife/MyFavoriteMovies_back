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

  repository.AddUser(c, &user)
  c.Status(http.StatusCreated)
}

func UpdateUserHandler(c *gin.Context) {
  user, errUser := utils.CheckContextUser(c)
  if !errUser || !utils.BindJSON(c, &user) {
      return
  }

  repository.UpdateUser(c, user)
}