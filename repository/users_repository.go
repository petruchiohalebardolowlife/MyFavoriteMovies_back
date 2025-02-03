package repository

import (
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
)

func AddUser(user *structs.User) error {
  if err := database.DB.Create(user).Error; err != nil {
      return err
  }
  return nil
}

func UpdateUser(user *structs.User) error {
  if err := database.DB.Save(user).Error; err != nil {
      return err
  }
  return nil
}