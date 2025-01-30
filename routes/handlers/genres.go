package handlers

import (
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AddGenres(c *gin.Context) {
    var input struct {
        Genres []struct {
            ID   uint   `json:"id"`
            Name string `json:"name"`
        } `json:"genres"`
    }

    if !utils.BindJSON(c, &input) {
        return
    }

    extractGenreIDs := func(genres []struct {
        ID   uint   `json:"id"`
        Name string `json:"name"`
    }) []uint {
        var ids []uint
        for _, genre := range genres {
            ids = append(ids, genre.ID)
        }
        return ids
    }

    var genres []structs.Genre
    if err := database.DB.Where("id IN (?)", extractGenreIDs(input.Genres)).Find(&genres).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    existingGenreMap := make(map[uint]struct{})
    for _, genre := range genres {
        existingGenreMap[genre.ID] = struct{}{}
    }

    var newGenres []structs.Genre
    for _, genreData := range input.Genres {
        if _, exists := existingGenreMap[genreData.ID]; exists {
            continue 
        }
        newGenres = append(newGenres, structs.Genre{
            ID:   genreData.ID,
            Name: genreData.Name,
        })
    }

    if len(newGenres) > 0 {
        if err := database.DB.Create(&newGenres).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    if len(newGenres) == 0 {
        c.JSON(http.StatusConflict, gin.H{"error": "All genres already exist"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Genres added successfully", "data": newGenres})
    fmt.Fprintln(os.Stdout, "GENRES ADDED")
}

func AddOrDeleteFavoriteGenre(c *gin.Context) {
    userID, ok, user := utils.CheckUser(c)
    if !ok {
        return
    }

    var input struct {
        GenreID uint `json:"genre_id"`
    }

    if !utils.BindJSON(c, &input) {
        return
    }

    var existingGenre structs.FavoriteGenre
    dbErr := database.DB.Where("user_id = ? AND genre_id = ?", userID, input.GenreID).First(&existingGenre).Error

    if dbErr == nil {
        if err := database.DB.Delete(&existingGenre).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Favorite genre removed successfully"})
        return
    }

        var genre structs.Genre
        if err := database.DB.First(&genre, input.GenreID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
            return
        }

        newFavorite := structs.FavoriteGenre{
            UserID:  uint(userID),
            GenreID: input.GenreID,
            User:    user,
            Genre:   genre,
        }

        if err := database.DB.Create(&newFavorite).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Favorite genre added successfully", "data": newFavorite})
}
