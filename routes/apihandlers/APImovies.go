package api

import (
	"encoding/json"
	"fmt"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
  user, errUser := utils.GetContextUser(c)
  if !errUser {
    return
}

  var favgenreIDs []string
  if err := database.DB.
        Model(&structs.FavoriteGenre{}).
        Where("user_id = ?", user.ID).
        Pluck("genre_id", &favgenreIDs).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No favorite genres found"})
        return
  }

  selectedGenre := c.Query("genre")
  minPopularity := c.Query("popularity")
  year := c.Query("year")    
  page := c.DefaultQuery("page", "1") 

  genreFilter := strings.Join(favgenreIDs, ",")
  if selectedGenre != "" {
        genreFilter = selectedGenre
  }

  apiKey := "6b2c0c7ec76b014687e6201bb7bd904d"
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
