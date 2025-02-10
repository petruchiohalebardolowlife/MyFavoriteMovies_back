package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
)

func AddUser(user *structs.User) error {
  if user.NickName == "" || user.UserName == "" || user.Password == "" {
		return errors.New("Some of filed are empty")
	}

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

func DeleteUser(userID int32) error {
  if err := database.DB.Where("id = ?", userID).Delete(&structs.User{}).Error; err != nil {
      return errors.New("User not found")
  }

  return nil
}