package api

import (
	"encoding/json"
	"fmt"
	config "myfavouritemovies/configs"
	"myfavouritemovies/structs"
	"net/http"
)

func FetchGenres() ([]structs.Genre, error) {
  url := fmt.Sprintf("%sgenre/movie/list?api_key=%s", config.TMDB_API_BASE_URL, config.API_KEY)
  resp,err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  var result struct {
    Genres []structs.Genre `json:"genres"`
  }

  if err:=json.NewDecoder(resp.Body).Decode(&result); err != nil {
    return nil, err
  }
  return result.Genres, nil
}