package service

import (
	"encoding/json"
	"myfavouritemovies/structs"

	"strconv"
)

func FetchMovieDetails(movie_id int) (structs.MovieDetails, error) {
	endpoint := "/movie/"+strconv.Itoa(movie_id)
	body, err := FetchFromTMDB(endpoint)
	if err != nil {
		return structs.MovieDetails{}, err
	}

	var response structs.MovieDetails

  if err := json.Unmarshal(body, &response); err != nil {
    return structs.MovieDetails{}, err
  }
  return response, nil
}