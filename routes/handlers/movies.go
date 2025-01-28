package handlers

import (
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddMovie(c *gin.Context) {
    var movie structs.Movie
    if err := c.ShouldBindJSON(&movie); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := database.DB.Create(&movie).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, movie)
	fmt.Fprintln(os.Stdout, "MOVIE ADD")
}

func AddFavouriteMovie (c *gin.Context) {
	userID,err:=strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid user ID"})
		return
	}

	var input struct {
		MovieID uint `json:"movie_id"`
	}
	if err := c.ShouldBindJSON(&input);err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
	return
	}
	var user structs.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

	var movie structs.Movie
	if err := database.DB.First(&movie, input.MovieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"Movie not found"})
		return
	}

	newFavourite := structs.FavouriteMovie {
    UserID: uint(userID),
    MovieID: input.MovieID,   
    Watched: false,
    User: user,
    Movie: movie,
	}

	if err := database.DB.Create(&newFavourite).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	c.JSON(http.StatusCreated, gin.H{"message": "Favorite genre added successfully", "data": newFavourite})
}

func ToggleWatchedStatus(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input struct {
		MovieID uint `json:"movie_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var favMovie structs.FavouriteMovie
	if err := database.DB.Where("user_id = ? AND movie_id = ?", userID, input.MovieID).First(&favMovie).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite movie not found"})
		return
	}

	favMovie.Watched = !favMovie.Watched

	if err := database.DB.Save(&favMovie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie watched status updated successfully", "data": favMovie})
}

func DeleteFavouriteMovie(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var input struct {
        MovieID uint `json:"movie_id"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existingMovie structs.FavouriteMovie
    if err := database.DB.Where("user_id = ? AND movie_id = ?", userID, input.MovieID).First(&existingMovie).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Favorite movie not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    if err := database.DB.Delete(&existingMovie).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Favorite movie deleted successfully"})
}