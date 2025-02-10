package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
)

func AddFavoriteMovie(userID int32, input structs.Movie) error {
  if _, err := utils.FindFavoriteMovie(userID, input.ID); err == nil {
      return errors.New("movie already in favorites")
  }

  newFavorite := structs.FavoriteMovie{
      UserID:      userID,
      MovieID:     (input.ID),
      Title:       input.Title,
      PosterPath:  input.PosterPath,
      VoteAverage: input.VoteAverage,
      Watched:     false,
  }

  if err := database.DB.Create(&newFavorite).Error; err != nil {
      return errors.New("failed to add favorite movie")
  }
  var genres []*structs.Genre
    if err := database.DB.Where("id IN ?", input.GenreIDs).Find(&genres).Error; err != nil {
        return errors.New("failed to find genres")
    }

  if err := database.DB.Model(&newFavorite).Association("Genres").Append(genres); err != nil {
      return errors.New("failed to associate genres")
  }

  return nil
}

func ToggleWatchedStatus(userID, movieID int32) error {
  favMovie, err := utils.FindFavoriteMovie(userID, movieID)
  if err != nil {
      return errors.New("favorite movie not found")
  }

  favMovie.Watched = !favMovie.Watched

  if err := database.DB.Save(&favMovie).Error; err != nil {
    return errors.New("failed to toggle")
  }

  return nil
}

func DeleteFavoriteMovie(userID, movieID int32) error {
  existingMovie, err := utils.FindFavoriteMovie(userID, movieID)
  if err != nil {
      return errors.New("favorite movie not found")
  }

  if err := database.DB.Delete(&existingMovie).Error; err != nil {
      return errors.New("failed to delete favorite movie")
  }

  return nil
}

func GetFavoriteMovies(userID int32) ([]*structs.FavoriteMovie, error) {
  var favMovies []*structs.FavoriteMovie
  if err := database.DB.
      Preload("Genres").
      Where("user_id = ?", userID).
      Find(&favMovies).Error; err != nil {
      return nil, errors.New("failed to get favorite movies")
  }
  
  return favMovies, nil
}