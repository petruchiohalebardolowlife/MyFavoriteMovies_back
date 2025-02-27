package service

import (
	"encoding/json"

	"myfavouritemovies/models"
)

func FetchGenres() ([]models.Genre, error) {
  endpoint := "/genre/movie/list"
  body, err := FetchFromTMDB(endpoint, "")
  if err != nil {
    return nil, err
  }

  var response struct {
    Genres []models.Genre `json:"genres"`
  }

  if err := json.Unmarshal(body, &response); err != nil {
    return nil, err
  }

  return response.Genres, nil
}
