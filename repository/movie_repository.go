package repository

import (
	"errors"
	"math"
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

func ToggleWatchedStatus(favMovieID, userID uint) (*models.FavoriteMovie, error) {
  favMovie, err := FindFavoriteMovie(favMovieID, userID)
  if err != nil {
    return nil, errors.New("favorite movie not found")
  }

  favMovie.Watched = !favMovie.Watched

  if err := database.DB.Save(&favMovie).Error; err != nil {
    return nil, errors.New("failed to toggle")
  }

  return &favMovie, nil
}

func DeleteFavoriteMovie(favMovieID, userID uint) error {
  existingMovie, err := FindFavoriteMovie(favMovieID, userID)
  if err != nil {
    return errors.New("favorite movie not found")
  }

  if err := database.DB.Delete(&existingMovie).Error; err != nil {
    return errors.New("failed to delete favorite movie")
  }

  return nil
}

func GetFavoriteMovies(userID uint, page, moviesPerPage uint) (*models.GetFavoriteMoviesResponse, error) {
  var favMovies []*models.FavoriteMovie
  var numberOfMovies int64

  if err := database.DB.Model(&models.FavoriteMovie{}).
      Where("user_id = ?", userID).
      Count(&numberOfMovies).Error; err != nil {
      return nil, errors.New("failed to count favorite movies")
  }

  totalPages := int(math.Ceil(float64(numberOfMovies) / float64(moviesPerPage)));
  if (page > uint(totalPages)) {
    return nil, errors.New("incorrect page")
  }
  offset := (page - 1) * moviesPerPage
  if err := database.DB.
      Preload("Genres").
      Where("user_id = ?", userID).
      Limit(int(moviesPerPage)).
      Offset(int(offset)).
      Find(&favMovies).Error; err != nil {
      return nil, errors.New("failed to get favorite movies")
  }

  return &models.GetFavoriteMoviesResponse{Page: uint(page), TotalPages: uint(totalPages), Results:favMovies}, nil
}

func FindFavoriteMovie(favMovieID, userID uint) (models.FavoriteMovie, error) {
  var favMovie models.FavoriteMovie
  if err := database.DB.Where("movie_id = ? AND user_id = ?", favMovieID, userID).First(&favMovie).Error; err != nil {
    return models.FavoriteMovie{}, err
  }

  return favMovie, nil
}
