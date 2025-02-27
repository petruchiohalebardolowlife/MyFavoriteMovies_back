package service

import (
	"fmt"
	"io"
	config "myfavouritemovies/configs"
	"net/http"
)

func FetchFromTMDB(endpoint string, params string) ([]byte, error) {
  url := fmt.Sprintf("%s%s?api_key=%s&%s", config.TMDB_API_BASE_URL, endpoint, config.API_KEY, params)

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
