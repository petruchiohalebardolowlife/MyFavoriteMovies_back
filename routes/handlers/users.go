package handlers

import (
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
    var user structs.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
	fmt.Fprintln(os.Stdout, "USER ADD!")
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user structs.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	fmt.Fprintln(os.Stdout, "USER UPDATED!")
}

func ReadUser(c *gin.Context) {
	id := c.Param("id")
	var user structs.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
	fmt.Fprintln(os.Stdout, "USER READ!")
}