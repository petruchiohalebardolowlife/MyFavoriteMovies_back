package handlers

import (
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
    var user structs.User
    if !utils.BindJSON(c, &user) {
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
	session := sessions.Default(c)
    userInterface := session.Get("user")
    user, errUser := userInterface.(structs.User)
    if !errUser {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    return
    }

	if !utils.BindJSON(c, &user) {
        return
    }
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	fmt.Fprintln(os.Stdout, "USER UPDATED!")
}
