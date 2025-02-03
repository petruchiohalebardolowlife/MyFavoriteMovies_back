func AddGenresOnStartup(db *gorm.DB) {
  apiKey := "6b2c0c7ec76b014687e6201bb7bd904d"
  url := fmt.Sprintf("https://api.themoviedb.org/3/genre/movie/list?api_key=%s&language=en-US", apiKey)

  resp, err := http.Get(url)
  if err != nil {
      log.Printf("Failed to fetch genres: %v", err)
      return
  }
  defer resp.Body.Close()

  var result struct {
      Genres []struct {
          ID   int    `json:"id"`
          Name string `json:"name"`
      } `json:"genres"`
  }

  if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      log.Printf("Failed to parse genres: %v", err)
      return
  }

  var addedGenres []structs.Genre

  for _, genreData := range result.Genres {
      var existingGenre structs.Genre
      err := db.Where("id = ?", genreData.ID).First(&existingGenre).Error
      if err != nil && err != gorm.ErrRecordNotFound {
          log.Printf("Error checking genre existence: %v", err)
          return
      }
      
      if err == nil {
          continue
      }

      newGenre := structs.Genre{
          ID:   uint(genreData.ID),
          Name: genreData.Name,
      }

      if err := db.Create(&newGenre).Error; err != nil {
          log.Printf("Failed to add genre to DB: %v", err)
          return
      }

      addedGenres = append(addedGenres, newGenre)
  }

  if len(addedGenres) == 0 {
      log.Println("All genres already exist")
  } else {
      log.Println("Genres added successfully")
  }
}
