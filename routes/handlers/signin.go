package handlers

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignIn(c *gin.Context) {
    var input structs.SignIn

    if !utils.BindJSON(c, &input) {
        return
    }

    var user structs.User
    if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err == gorm.ErrRecordNotFound {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    if user.Password != input.Password {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or Password"})
        return
    }

    session := sessions.Default(c)
    session.Set("user", user)
    if err := session.Save(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save session"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}