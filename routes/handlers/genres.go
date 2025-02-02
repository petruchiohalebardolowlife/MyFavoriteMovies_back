package handlers

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func AddGenres(c *gin.Context) {
	var input struct {
		Genres []structs.Genre `json:"genres"`
	}

    if !utils.BindJSON(c, &input) {
        return
    }

	var existingGenres []structs.Genre
	if err := database.DB.Find(&existingGenres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var newGenres []structs.Genre
	for _, genre := range input.Genres {
		if !slices.ContainsFunc(existingGenres, func(gen structs.Genre) bool {
			return gen.ID == genre.ID
		}) {
			newGenres = append(newGenres, genre)
		}
	}

	if len(newGenres) > 0 {
		if err := database.DB.Create(&newGenres).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Genres added successfully"})
	} else {
		c.Status(http.StatusCreated)
	}
}


func AddFavoriteGenre(c *gin.Context) {
    var input struct {
        GenreID uint `json:"genre_id"`
    }

    user, errUser := utils.CheckSession(c)
    if !errUser || !utils.BindJSON(c, &input) {
        return
    }

    dbErr := database.DB.Where("user_id = ? AND genre_id = ?", user.ID, input.GenreID).First(&structs.FavoriteGenre{}).Error
    if dbErr == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Genre already in favorites"})
        return
    }

    newFavorite := structs.FavoriteGenre{
        UserID:  user.ID,
        GenreID: input.GenreID,
    }

    if err := database.DB.Create(&newFavorite).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Favorite genre added successfully"})
}

func DeleteFavoriteGenre(c *gin.Context) {
    
    var input struct {
        GenreID uint `json:"genre_id"`
    }
    user, errUser := utils.CheckSession(c)
    if !errUser || !utils.BindJSON(c, &input) {
        return
    }

    var favGenre structs.FavoriteGenre
    if err := database.DB.Where("user_id = ? AND genre_id = ?", user.ID, input.GenreID).First(&favGenre).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found in favorites"})
    return
    }

    if err := database.DB.Delete(&favGenre).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Favorite genre deleted successfully"})
}

