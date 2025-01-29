package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddGenresOnStartup(db *gorm.DB) {
    apiKey := "6b2c0c7ec76b014687e6201bb7bd904d"
    url := fmt.Sprintf("https://api.themoviedb.org/3/genre/movie/list?api_key=%s&language=en-US", apiKey)

    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Failed to fetch genres: %v", err)
        return
    }
    defer resp.Body.Close()

    var result struct {
        Genres []struct {
            ID   int    `json:"id"`
            Name string `json:"name"`
        } `json:"genres"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Printf("Failed to parse genres: %v", err)
        return
    }

    var addedGenres []structs.Genre

    for _, genreData := range result.Genres {
        var existingGenre structs.Genre
        err := db.Where("id = ?", genreData.ID).First(&existingGenre).Error
        if err != nil && err != gorm.ErrRecordNotFound {
            log.Printf("Error checking genre existence: %v", err)
            return
        }
        
        if err == nil {
            continue
        }

        newGenre := structs.Genre{
            ID:   uint(genreData.ID),
            Name: genreData.Name,
        }

        if err := db.Create(&newGenre).Error; err != nil {
            log.Printf("Failed to add genre to DB: %v", err)
            return
        }

        addedGenres = append(addedGenres, newGenre)
    }

    if len(addedGenres) == 0 {
        log.Println("All genres already exist")
    } else {
        log.Println("Genres added successfully")
    }
}


func AddFavoriteGenres(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var input struct {
        GenreIDs []uint `json:"genre_ids"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user structs.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    var addedGenres []structs.FavoriteGenre

    for _, genreID := range input.GenreIDs {
        var existingGenre structs.FavoriteGenre
        err := database.DB.Where("user_id = ? AND genre_id = ?", userID, genreID).First(&existingGenre).Error

        if err == nil {
            continue
        }

        if err != gorm.ErrRecordNotFound {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var genre structs.Genre
        if err := database.DB.First(&genre, genreID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "One or more genres not found"})
            return
        }

        newFavourite := structs.FavoriteGenre{
            UserID:  uint(userID),
            GenreID: genreID,
            User:    user,
            Genre:   genre,
        }

        if err := database.DB.Create(&newFavourite).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        addedGenres = append(addedGenres, newFavourite)
    }

    if len(addedGenres) == 0 {
        c.JSON(http.StatusConflict, gin.H{"error": "All genres are already marked as favorite"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Favorite genres added successfully", "data": addedGenres})
}

func DeleteFavoriteGenre(c *gin.Context) {
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

    var existingGenre structs.FavoriteGenre
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

