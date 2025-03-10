package service

import (
	"encoding/json"
	"myfavouritemovies/models"

	"strconv"
)

func FetchMovieDetails(movieID uint) (*models.MovieDetails, error) {
  endpoint := "/movie/" + strconv.Itoa(int(movieID))
  body, err := FetchFromTMDB(endpoint, "")
  if err != nil {
    return nil, err
  }

  var response models.MovieDetails

  if err := json.Unmarshal(body, &response); err != nil {
    return nil, err
  }

  return &response, nil
}
