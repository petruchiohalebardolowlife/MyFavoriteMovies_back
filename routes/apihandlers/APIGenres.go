package apihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAPIGenres(c *gin.Context) {
    apiKey := "6b2c0c7ec76b014687e6201bb7bd904d"
    url := fmt.Sprintf("https://api.themoviedb.org/3/genre/movie/list?api_key=%s&language=en-US", apiKey)

    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genres"})
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse genres"})
        return
    }

    c.JSON(http.StatusOK, result.Genres)
}
