package routes

import (
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.POST("/movies",addMovie)
	router.POST("/users",addUser)
}

func addMovie(c *gin.Context) {
    var movie models.Movie
    if err := c.ShouldBindJSON(&movie); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := database.DB.Create(&movie).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, movie)
	fmt.Fprintln(os.Stdout, "MOVIEADDED")
}

func addUser(c *gin.Context) {
    var user models.User
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