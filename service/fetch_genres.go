package service

import (
	"encoding/json"

	"myfavouritemovies/structs"
)

func FetchGenres() ([]structs.Genre, error) {
  endpoint := "/genre/movie/list"
  body, err := FetchFromTMDB(endpoint,"")
  if err != nil {
    return nil, err
  }

  var response struct {
    Genres []structs.Genre `json:"genres"`
  }

  if err := json.Unmarshal(body, &response); err != nil {
    return nil, err
  }
  return response.Genres, nil
}