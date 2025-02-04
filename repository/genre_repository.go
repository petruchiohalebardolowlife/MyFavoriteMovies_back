package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func GetAllGenres() ([]structs.Genre, error) {
	var genres []structs.Genre
	if err := database.DB.Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

func SaveGenresToDB(db *gorm.DB, genres []structs.Genre) error {
	existingGenres, err := GetAllGenres()
  if err != nil {
      return err
  }

	var newGenres []structs.Genre
	for _, genre := range genres {
		if !slices.ContainsFunc(existingGenres, func(gen structs.Genre) bool {
			return gen.ID == genre.ID
		}) {
			newGenres = append(newGenres, genre)
		}
	}

	if len(newGenres) > 0 {
		if err := db.Create(&newGenres).Error; err != nil {
			return err
		}
	}

	return nil
}

func AddFavoriteGenre(userID, genreID uint) error {
  err := database.DB.Where("user_id = ? AND genre_id = ?", userID, genreID).
      First(&structs.FavoriteGenre{}).Error
  if err == nil {
      return errors.New("genre already in favorites")
  }

  newFavorite := structs.FavoriteGenre{
      UserID:  userID,
      GenreID: genreID,
  }

  if err := database.DB.Create(&newFavorite).Error; err != nil {
      return err
  }

  return nil
}

func DeleteFavoriteGenre(userID, genreID uint) error {
  var favGenre structs.FavoriteGenre
  if err := database.DB.Where("user_id = ? AND genre_id = ?", userID, genreID).
      First(&favGenre).Error; err != nil {
      return errors.New("genre not in favorites")
  }

  if err := database.DB.Delete(&favGenre).Error; err != nil {
      return err
  }
  
  return nil
}