package tmdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	config "myfavouritemovies/configs"
	"myfavouritemovies/structs"
	"net/http"
)

func FetchFromTMDB(endpoint string) ([]byte, error) {
  url := fmt.Sprintf("%s%s?apikey=%s", config.TMDB_API_BASE_URL, endpoint, config.API_KEY)
  fmt.Println(url)
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

	// Используем json.NewDecoder для потокового чтения
	decoder := json.NewDecoder(bytes.NewReader(body))

	// Декодируем данные
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	// Проверяем, что жанры есть в ответе
	if len(response.Genres) == 0 {
		return nil, fmt.Errorf("no genres found in the response")
	}

	return response.Genres, nil
}
