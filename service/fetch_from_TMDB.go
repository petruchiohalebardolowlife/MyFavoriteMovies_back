package service

import (
	"fmt"
	"io"
	config "myfavouritemovies/configs"
	"net/http"
	"strings"
)

func FetchFromTMDB(endpoint string) ([]byte, error) {
  separator := "?"
  if strings.Contains(endpoint, "?") {
    separator = "&"
  }

  url := fmt.Sprintf("%s%s%sapi_key=%s", config.TMDB_API_BASE_URL, endpoint, separator, config.API_KEY)

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


