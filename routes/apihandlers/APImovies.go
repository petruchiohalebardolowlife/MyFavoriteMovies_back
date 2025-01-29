package apihandlers

import (
	"encoding/json"
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var favouriteGenres []structs.FavouriteGenre
    if err := database.DB.Where("user_id = ?", userID).Find(&favouriteGenres).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No favorite genres found"})
        return
    }

    allGenreIDs := []string{}
    for _, genre := range favouriteGenres {
        allGenreIDs = append(allGenreIDs, fmt.Sprintf("%d", genre.GenreID))
    }

    selectedGenre := c.Query("genre")
    minPopularity := c.Query("popularity")
    year := c.Query("year")    
    page := c.DefaultQuery("page", "1") 

    genreFilter := strings.Join(allGenreIDs, ",")
    if selectedGenre != "" {
        genreFilter = selectedGenre
    }

    apiKey := "YOUR_TMDB_API_KEY"
    url := fmt.Sprintf(
        "https://api.themoviedb.org/3/discover/movie?api_key=%s&with_genres=%s&page=%s",
        apiKey, genreFilter, page,
    )

    if minPopularity != "" {
        url += fmt.Sprintf("&vote_average.gte=%s", minPopularity)
    }
    if year != "" {
        url += fmt.Sprintf("&primary_release_year=%s", year)
    }

    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
        return
    }
    defer resp.Body.Close()

    var result struct {
        Results []struct {
            ID          int     `json:"id"`
            Title       string  `json:"title"`
            Popularity  float64 `json:"popularity"`
            ReleaseDate string  `json:"release_date"`
        } `json:"results"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse movies"})
        return
    }

    c.JSON(http.StatusOK, result.Results)
}
