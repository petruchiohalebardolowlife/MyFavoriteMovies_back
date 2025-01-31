package handlers

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)


func AddFavoriteMovie(c *gin.Context) {
    var input structs.Movie

    user, errUser := utils.CheckSession(c)
    if !errUser || !utils.BindJSON(c, &input) {
        return
    }

    if _, err := utils.FindFavoriteMovie(user.ID, input.MovieID); err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Movie already in favorites"})
        return
    }

    newFavorite := structs.FavoriteMovie{
        UserID:     user.ID,
        MovieID:    input.MovieID,
        Title:      input.Title,
        PosterPath: input.PosterPath,
        VoteAverage: input.VoteAverage,
        Genres: input.Genres,
        Watched:    false,
    }

    if err := database.DB.Create(&newFavorite).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if err := database.DB.Model(&newFavorite).Association("Genres").Append(input.Genres); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate genres"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Movie added to favorites successfully", "data": newFavorite})
}

func ToggleWatchedStatus(c *gin.Context) {
    session := sessions.Default(c)
    userInterface := session.Get("user")
    user, errUser := userInterface.(structs.User)
    if !errUser {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    return
    }
	var input struct {
		MovieID uint `json:"movie_id"`
	}
    if !utils.BindJSON(c, &input) {
        return
    }

    favMovie, err := utils.FindFavoriteMovie(uint(user.ID), input.MovieID)
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
    session := sessions.Default(c)
    userInterface := session.Get("user")
    user, errUser := userInterface.(structs.User)
    if !errUser {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    return
    }

    var input struct {
        MovieID uint `json:"movie_id"`
    }
    if !utils.BindJSON(c, &input) {
        return
    }

    existingMovie, err := utils.FindFavoriteMovie(uint(user.ID), input.MovieID)
    if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find favorite movie"})
    return
    }

    database.DB.Delete(&existingMovie)
    c.JSON(http.StatusOK, gin.H{"message": "Favorite movie deleted successfully"})
}   