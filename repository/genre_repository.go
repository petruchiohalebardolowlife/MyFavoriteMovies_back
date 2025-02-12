package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"

	"golang.org/x/exp/slices"
)

func GetAllGenres() ([]*models.Genre, error) {
  var genres []*models.Genre
  if err := database.DB.Find(&genres).Error; err != nil {
    return nil, err
  }

  return genres, nil
}

func SaveGenresToDB(genres []models.Genre) error {
  existingGenres, err := GetAllGenres()
  if err != nil {
      return err
  }

  var newGenres []models.Genre
  for _, genre := range genres {
    if !slices.ContainsFunc(existingGenres, func(gen *models.Genre) bool {
      return gen.ID == genre.ID
    }) {
      newGenres = append(newGenres, genre)
    }
  }

  if len(newGenres) > 0 {
    if err := database.DB.Create(&newGenres).Error; err != nil {
      return err
    }
  }

  return nil
}

func AddFavoriteGenre(userID, genreID uint) error {
  err := database.DB.Where("user_id = ? AND genre_id = ?", userID, genreID).
      First(&models.FavoriteGenre{}).Error
  if err == nil {
      return errors.New("genre already in favorites")
  }

  newFavorite := models.FavoriteGenre{
      UserID:  userID,
      GenreID: genreID,
  }

  if err := database.DB.Create(&newFavorite).Error; err != nil {
      return err
  }

  return nil
}

func DeleteFavoriteGenre(userID, genreID uint) error {
  var favGenre models.FavoriteGenre
  if err := database.DB.Where("user_id = ? AND genre_id = ?", userID, genreID).
      First(&favGenre).Error; err != nil {
      return errors.New("genre not in favorites")
  }

  if err := database.DB.Delete(&favGenre).Error; err != nil {
      return err
  }
  
  return nil
}

func GetFavoriteGenres (userID uint) ([]uint ,error) {
  var favGenresIDs []uint
  if err := database.DB.Model(&models.FavoriteGenre{}).Where("user_id = ?", userID).Pluck("genre_id", &favGenresIDs).Error; err != nil {
    return nil, errors.New("No favorite genres found")
}
  
  return favGenresIDs, nil
}
