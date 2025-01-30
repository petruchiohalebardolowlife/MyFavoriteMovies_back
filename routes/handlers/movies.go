package handlers

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func AddFavoriteMovie(c *gin.Context) {
    userID, ok, user := utils.CheckUser(c)
    if !ok {
        return
    }

    var input struct {
        MovieID     uint    `json:"movie_id"`
        Title       string  `json:"title"`
        PosterPath  string  `json:"poster_path"`
        VoteAverage float64 `json:"vote_average"`
    }

    if !utils.BindJSON(c, &input) {
        return
    }

    if _, err := utils.FindFavoriteMovie(uint(userID), input.MovieID); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Movie already in favorites"})
		return
	}
    
    newFavorite := structs.FavoriteMovie{
        UserID:     uint(userID),
        MovieID:    input.MovieID,
        Title:      input.Title,
        PosterPath: input.PosterPath,
        VoteAverage: input.VoteAverage,
        Watched:    false,
        User:       user,
    }

    if err := database.DB.Create(&newFavorite).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Movie added to favorites successfully", "data": newFavorite})
}

func ToggleWatchedStatus(c *gin.Context) {
    userID, ok, _ := utils.CheckUser(c)
    if !ok {
        return
    }
	var input struct {
		MovieID uint `json:"movie_id"`
	}
    if !utils.BindJSON(c, &input) {
        return
    }

    favMovie, err := utils.FindFavoriteMovie(uint(userID), input.MovieID)
	if err != nil {
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

func DeleteFavoriteMovie(c *gin.Context) {
    userID, ok, _ := utils.CheckUser(c)
    if !ok {
        return
    }

    var input struct {
        MovieID uint `json:"movie_id"`
    }
    if !utils.BindJSON(c, &input) {
        return
    }

    existingMovie, err := utils.FindFavoriteMovie(uint(userID), input.MovieID)
    if err != nil {
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