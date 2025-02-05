package service

import (
	"encoding/json"
	"fmt"
	"myfavouritemovies/structs"

	"strconv"
	"strings"
)

func FetchFiltredMovies(filters structs.MovieFilter) ([]structs.Movie, error){
  endpoint := "/discover/movie"

  var queryParams []string
  if len(filters.GenreIDs) > 0 {
    var genreStrings []string
    for _, id := range filters.GenreIDs {
      genreStrings = append(genreStrings, strconv.Itoa(id))
    }
    queryParams = append(queryParams, fmt.Sprintf("with_genres=%s", strings.Join(genreStrings, ",")))
  }
	if filters.Popularity > 0 {
		queryParams = append(queryParams, fmt.Sprintf("vote_average.gte=%f", filters.Popularity))
	}
	if filters.Year > 0 {
		queryParams = append(queryParams, fmt.Sprintf("primary_release_year=%d", filters.Year))
	}
	if filters.Page > 0 {
		queryParams = append(queryParams, fmt.Sprintf("page=%d", filters.Page))
	}
  var paramsString string
	if len(queryParams) > 0 {
		 paramsString = fmt.Sprintf("%s?%s", endpoint, strings.Join(queryParams, "&"))
	}

	body, err := FetchFromTMDB(endpoint, paramsString)
	if err != nil {
		return nil, err
	}

	var response structs.ResponseFiltredMovies
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
  
	return response.Results, nil
}