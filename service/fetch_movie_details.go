package service

import (
	"encoding/json"
	"myfavouritemovies/structs"

	"strconv"
)

func FetchMovieDetails(movieID uint) (*structs.MovieDetails, error) {
  endpoint := "/movie/"+strconv.Itoa(int(movieID))
  body, err := FetchFromTMDB(endpoint,"")
  if err != nil {
    return nil, err
  }

  var response structs.MovieDetails

  if err := json.Unmarshal(body, &response); err != nil {
    return nil, err
  }
  
  return &response, nil
}