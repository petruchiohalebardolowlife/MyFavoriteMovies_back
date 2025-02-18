package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
)

func AddFavoriteMovie(userID uint, input models.MovieInput) (*models.FavoriteMovie, error) {
  if err := database.DB.Where("user_id = ? AND movie_id = ?", userID, input.MovieID).First(&models.FavoriteMovie{}).Error; err == nil {
      return nil, errors.New("movie already in favorites")
  }

  newFavorite := models.FavoriteMovie{
      UserID:      userID,
      MovieID:     input.MovieID,
      Title:       input.Title,
      PosterPath:  input.PosterPath,
      VoteAverage: input.VoteAverage,
      Watched:     false,
  }

  if err := database.DB.Create(&newFavorite).Error; err != nil {
      return nil, errors.New("failed to add favorite movie")
  }
  var genres []*models.Genre
    if err := database.DB.Where("id IN ?", input.GenreIDs).Find(&genres).Error; err != nil {
        return nil, errors.New("failed to find genres")
    }

  if err := database.DB.Model(&newFavorite).Association("Genres").Append(genres); err != nil {
      return nil, errors.New("failed to associate genres")
  }

  return &newFavorite, nil
}

func ToggleWatchedStatus(favMovieID uint) (*models.FavoriteMovie, error) {
  favMovie, err := FindFavoriteMovie(favMovieID)
  if err != nil {
      return nil, errors.New("favorite movie not found")
  }

  favMovie.Watched = !favMovie.Watched

  if err := database.DB.Save(&favMovie).Error; err != nil {
    return nil, errors.New("failed to toggle")
  }

  return &favMovie, nil
}

func DeleteFavoriteMovie(favMovieID uint) error {
  existingMovie, err := FindFavoriteMovie(favMovieID)
  if err != nil {
      return errors.New("favorite movie not found")
  }

  if err := database.DB.Delete(&existingMovie).Error; err != nil {
      return errors.New("failed to delete favorite movie")
  }

  return nil
}

func GetFavoriteMovies(userID uint) ([]*models.FavoriteMovie, error) {
  var favMovies []*models.FavoriteMovie
  if err := database.DB.
      Preload("Genres").
      Where("user_id = ?", userID).
      Find(&favMovies).Error; err != nil {
      return nil, errors.New("failed to get favorite movies")
  }
  
  return favMovies, nil
}

func FindFavoriteMovie(favMovieID uint) (models.FavoriteMovie, error) {
  var favMovie models.FavoriteMovie
  if err := database.DB.Where("id = ?", favMovieID).First(&favMovie).Error; err != nil{
    return models.FavoriteMovie{},err
  }
  
  return favMovie, nil
}

