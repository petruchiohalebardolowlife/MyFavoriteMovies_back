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

func AddGenre(c *gin.Context) {
    var genre structs.Genre
    if err := c.ShouldBindJSON(&genre); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := database.DB.Create(&genre).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, genre)
    fmt.Fprintln(os.Stdout, "GENRE ADDED")
}

func AddFavouriteGenre(c *gin.Context) {

    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var input struct {
        GenreID uint `json:"genre_id"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existingGenre structs.FavouriteGenre
    if err := database.DB.Where("user_id = ? AND genre_id = ?", userID, input.GenreID).First(&existingGenre).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "This genre is already marked as favorite."})
        return
    }

    var user structs.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    var genre structs.Genre
    if err := database.DB.First(&genre, input.GenreID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
        return
    }


    newFavourite := structs.FavouriteGenre{
        UserID:  uint(userID),
        GenreID: input.GenreID,
        User:    user,
        Genre:   genre,
    }

    if err := database.DB.Create(&newFavourite).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Favorite genre added successfully", "data": newFavourite})
}

func DeleteFavouriteGenre(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var input struct {
        GenreID uint `json:"genre_id"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existingGenre structs.FavouriteGenre
    if err := database.DB.Where("user_id = ? AND genre_id = ?", userID, input.GenreID).First(&existingGenre).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Favorite genre not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    if err := database.DB.Delete(&existingGenre).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Favorite genre deleted successfully"})
}

