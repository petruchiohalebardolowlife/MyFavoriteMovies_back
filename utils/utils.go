package utils

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindJSON (c *gin.Context, input interface{}) bool {
	if err:=c.ShouldBindJSON(input);err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return false
	}
	return true
}

func CheckUser (c *gin.Context) (int, bool, structs.User) {
	userID,errID :=strconv.Atoi(c.Param("id"))
	if errID != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "Invalid user ID"})
		return 0, false, structs.User{}
	}
	
	var user structs.User
	if errDB := database.DB.First(&user, userID).Error; errDB != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return 0, false, structs.User{}
	}
	return userID, true, user
}

func FindFavoriteMovie(userID, movieID uint) (structs.FavoriteMovie, error) {
	var favMovie structs.FavoriteMovie
	err := database.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&favMovie).Error
	return favMovie, err
}