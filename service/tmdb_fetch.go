package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	config "myfavouritemovies/configs"
	"myfavouritemovies/structs"
	"net/http"
	"strconv"
	"strings"
)

func FetchFromTMDB(endpoint string) ([]byte, error) {
  var url string
  if strings.Contains(endpoint, "?") {
    url = fmt.Sprintf("%s%s&api_key=%s", config.TMDB_API_BASE_URL, endpoint, config.API_KEY)
  } else {
    url = fmt.Sprintf("%s%s?api_key=%s", config.TMDB_API_BASE_URL, endpoint, config.API_KEY)
  }

  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  return body, nil
}

func FetchGenres() ([]structs.Genre, error) {
	endpoint := "/genre/movie/list"
	body, err := FetchFromTMDB(endpoint)
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
	if filters.Rating != "" {
		queryParams = append(queryParams, fmt.Sprintf("vote_average.gte=%s", filters.Rating))
	}
	if filters.Year > 0 {
		queryParams = append(queryParams, fmt.Sprintf("primary_release_year=%d", filters.Year))
	}
	if filters.Page > 0 {
		queryParams = append(queryParams, fmt.Sprintf("page=%d", filters.Page))
	}

	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, strings.Join(queryParams, "&"))
	}

	body, err := FetchFromTMDB(endpoint)
	if err != nil {
		return nil, err
	}

	var response struct {
    Page int `json:"page"`
		Results []structs.Movie `json:"results"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
for _,n := range response.Results {
  log.Printf("%s", n.Title)
}
  log.Printf("LOG FROM FILTREDMOVIE ENDPOINT IS %s",endpoint)
	return response.Results, nil
}

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
