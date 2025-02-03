package utils

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"

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

func HardcodedUserMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        user := structs.User{
            ID:       671,
            NickName: "Vasiliy",
            Username: "Vasya",
        }

        c.Set("user", user)
        c.Next()
    }
}

func GetContextUser(c *gin.Context) (*structs.User, bool) {
    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return nil, false
    }

    user, errUser := userInterface.(structs.User)
    if !errUser {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user from context"})
        return nil, false
    }

    return &user, true
}
