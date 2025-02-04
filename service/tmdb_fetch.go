package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	config "myfavouritemovies/configs"
	"myfavouritemovies/structs"
	"net/http"
)

func FetchFromTMDB(endpoint string) ([]byte, error) {
  url := fmt.Sprintf("%s%s?api_key=%s", config.TMDB_API_BASE_URL, endpoint, config.API_KEY)
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
