package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
	"myfavouritemovies/utils"
)

func AddFavoriteMovie(userID uint, input structs.Movie) error {
  if _, err := utils.FindFavoriteMovie(userID, uint(input.ID)); err == nil {
      return errors.New("movie already in favorites")
  }

  newFavorite := structs.FavoriteMovie{
      UserID:      userID,
      MovieID:     uint(input.ID),
      Title:       input.Title,
      PosterPath:  input.PosterPath,
      VoteAverage: input.VoteAverage,
      Genres:      input.Genres,
      Watched:     false,
  }

  if err := database.DB.Create(&newFavorite).Error; err != nil {
      return errors.New("failed to add favorite movie")
  }

  if err := database.DB.Model(&newFavorite).Association("Genres").Append(input.Genres); err != nil {
      return errors.New("failed to associate genres")
  }

  return nil
}

func ToggleWatchedStatus(userID, movieID uint) error {
  favMovie, err := utils.FindFavoriteMovie(userID, movieID)
  if err != nil {
      return errors.New("favorite movie not found")
  }

  favMovie.Watched = !favMovie.Watched

  if err := database.DB.Save(&favMovie).Error; err != nil {
      return err
  }

  return nil
}

func DeleteFavoriteMovie(userID, movieID uint) error {
  existingMovie, err := utils.FindFavoriteMovie(userID, movieID)
  if err != nil {
      return errors.New("favorite movie not found")
  }

  if err := database.DB.Delete(&existingMovie).Error; err != nil {
      return err
  }

  return nil
}