package service

import (
	"encoding/json"
	"myfavouritemovies/structs"

	"strconv"
)

func FetchMovieDetails(movieID int32) (*structs.MovieDetails, error) {
  endpoint := "/movie/"+strconv.FormatInt(int64(movieID), 10)
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