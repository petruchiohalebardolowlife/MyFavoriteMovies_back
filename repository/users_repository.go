package repository

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context, user *structs.User) {
  if err := database.DB.Create(user).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
}

func UpdateUser(c *gin.Context, user *structs.User) {
  if err := database.DB.Save(user).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
  c.Status(http.StatusOK)
}