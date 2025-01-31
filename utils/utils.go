package utils

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func BindJSON (c *gin.Context, input interface{}) bool {
	if err:=c.ShouldBindJSON(input);err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return false
	}
	return true
}


func FindFavoriteMovie(userID, movieID uint) (structs.FavoriteMovie, error) {
	var favMovie structs.FavoriteMovie
	err := database.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&favMovie).Error
	return favMovie, err
}

func CheckSession(c *gin.Context) (*structs.User, bool) {
	session := sessions.Default(c)
	userInterface := session.Get("user")

	user, ok := userInterface.(structs.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil, false
	}

	return &user, true
}